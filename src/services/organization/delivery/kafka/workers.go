package kafka

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"file_reader/src/instrument"
	"file_reader/src/models"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/avast/retry-go"
	"github.com/riferrei/srclient"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

const (
	retryAttempts = 1
	retryDelay    = 1 * time.Second
)

func (ocg *OrganizationsConsumerGroup) createOrganizationWorker(
	ctx context.Context,
	cancel context.CancelFunc,
	r *kafka.Reader,
	w *kafka.Writer,
	wg *sync.WaitGroup,
	workerID int,

) {
	defer wg.Done()
	defer cancel()

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			ocg.log.Error("Fetch message", zap.String("error", err.Error()))
			return
		}

		ocg.log.Sugar().Infof(

			"WORKER: %v, message at topic/parition/offset",
			zap.String("workerID", strconv.Itoa(workerID)),
			zap.String("Topic", m.Topic),
			zap.String("Partition", strconv.Itoa(m.Partition)),
			zap.String("Offset", string(m.Offset)),
			zap.String("Key", string(m.Key)),
			zap.String("Value", string(m.Value)),
		)

		var org models.Organization
		// This code is for testing only
		schemaID := binary.BigEndian.Uint32(m.Value[1:5])
		srcUrl := instrument.MustGetEnv("SCHEMA_REGISTRY_URL")
		schemaRegistryClient := srclient.CreateSchemaRegistryClient(srcUrl)
		schema, err := schemaRegistryClient.GetSchema(int(schemaID))
		if err != nil {
			panic(fmt.Sprintf("Error getting the schema with id '%d' %s", schemaID, err))
		}
		native, _, _ := schema.Codec().NativeFromBinary(m.Value[5:])
		value, _ := schema.Codec().TextualFromNative(nil, native)

		if err := json.Unmarshal(value, &org); err != nil {
			ocg.log.Sugar().Errorf("json.Unmarshal", err)
			continue
		}

		if err := retry.Do(func() error {
			created, err := ocg.organizationUC.Create(ctx, &org)
			if err != nil {
				return err
			}
			ocg.log.Sugar().Infof("created organization: %v", created)
			return nil
		},
			retry.Attempts(retryAttempts),
			retry.Delay(retryDelay),
			retry.Context(ctx),
		); err != nil {
			if err := ocg.publishErrorMessage(ctx, w, m, err); err != nil {
				ocg.log.Sugar().Errorf("publishErrorMessage", err)
				continue
			}
			ocg.log.Sugar().Errorf("organizationUC.Create.publishErrorMessage", err)
			continue
		}

		if err := r.CommitMessages(ctx, m); err != nil {
			ocg.log.Sugar().Errorf("CommitMessages", err)
			continue
		}

		// inc successMessages
	}
}
