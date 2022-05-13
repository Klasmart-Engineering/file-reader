package proto

import (
	"context"
	"encoding/csv"
	"file_reader/src/instrument"
	"file_reader/src/log"
	orgPb "file_reader/src/protos/onboarding"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"file_reader/src/third_party/protobuf"

	csvError "file_reader/src/pkg"

	"github.com/riferrei/srclient"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
)

type Config struct {
	BrokerAddrs []string
	Reader      io.Reader
	Context     context.Context
	Logger      *log.ZapLogger
}

type Operation struct {
	topic            string
	schema           *srclient.Schema
	rowToProtoSchema func(row []string, trackingId string) (*orgPb.Organization, error)
}

func (op Operation) GetNewKafkaWriter(config Config) *kafka.Writer {

	writerRequiredAcks, _ := strconv.Atoi(instrument.MustGetEnv("WRITE_REQUIRED_ACKS"))
	writerMaxAttempts, _ := strconv.Atoi(instrument.MustGetEnv("WRITE_MAX_ATTEMPTS"))
	writerReadTimeout, _ := strconv.Atoi(instrument.MustGetEnv("WRITE_READ_TIMEOUT"))

	writerWriteTimeout, _ := strconv.Atoi(instrument.MustGetEnv("WRITE_WRITE_TIMEOUT"))
	w := &kafka.Writer{
		Addr:         kafka.TCP(config.BrokerAddrs...),
		Topic:        op.topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequiredAcks(writerRequiredAcks),
		MaxAttempts:  writerMaxAttempts,
		//Logger:       kafka.LoggerFunc(config.Logger.),
		Compression:  compress.Snappy,
		ReadTimeout:  time.Duration(writerReadTimeout * int(time.Second)),
		WriteTimeout: time.Duration(writerWriteTimeout * int(time.Second)),
	}
	return w
}

func (op Operation) IngestFilePROTO(config Config, fileTypeName string, trackingId string) error {

	switch fileTypeName {

	case "CSV":
		csvReader := csv.NewReader(config.Reader)
		w := op.GetNewKafkaWriter(config)
		schemaID, err := schemaRegistryClient.GetProtoSchemaID(organizationSchemaName, organizationProtoTopic)
		if err != nil {
			return err
		}
		serde := protobuf.NewProtoSerDe()

		fails := 0 // Keep track of how many rows failed to process
		total := 0 // Total of row
		for {
			row, err := csvReader.Read()
			total += 1
			if err == io.EOF {
				break
			}
			if err != nil {
				config.Logger.Errorf(config.Context, false, err.Error())
				// If there is an error when reading the current row then skip it
				err = errors.Wrap(csvError.CSVRowError{RowNum: total, What: "Can't read the data"}, fmt.Sprintf("Error: %w", err))
				fails += 1
				continue
			}
			// Process the row, if it fails then skip that row
			orgSchema, err := op.rowToProtoSchema(row, trackingId)
			if err != nil {
				config.Logger.Errorf(config.Context, false, err.Error())
				err = errors.Wrap(csvError.CSVRowError{RowNum: total, What: "Fail to process"}, fmt.Sprintf("Error: %w", err))
				fails += 1
				continue
			}

			valueBytes, err := serde.Serialize(schemaID, orgSchema)

			if err != nil {
				config.Logger.Errorf(config.Context, false, fmt.Sprintf("error serializing message: %w", err))
				err = errors.Wrap(err, fmt.Sprintf("error serializing message: %w", err))
				return err
			}

			// Put the row on the topic
			err = w.WriteMessages(
				config.Context,
				kafka.Message{
					Key:   []byte(""),
					Value: valueBytes,
				},
			)
			if err != nil {
				config.Logger.Errorf(config.Context, false, fmt.Sprintf("could not write message: %w", err))
				err = errors.Wrap(err, fmt.Sprintf("could not write message: %w", err))
				return err
			}
		}
		// If all the rows are failed to process then
		if total == fails {
			config.Logger.Errorf(config.Context, false, err.Error())
			err = errors.Wrap(err, "All rows are invalid")
			return err
		}
	}
	return nil
}
