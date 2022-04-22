package usecase

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"file_reader/src/models"
	orgKafka "file_reader/src/services/organization/delivery/kafka"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/riferrei/srclient"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type organizationUC struct {
	log         zap.Logger
	orgProducer orgKafka.OrganizationProducer
}

// NewOrganizationUC contructor

func NewOrganizationUC(
	log zap.Logger, orgProducer orgKafka.OrganizationProducer) *organizationUC {
	return &organizationUC{log: log, orgProducer: orgProducer}
}

// Create new organization
func (o *organizationUC) Create(ctx context.Context, org *models.Organization) (*models.Organization, error) {
	return org, nil
}

func (o *organizationUC) PublishCreate(ctx context.Context, org *models.Organization, schema *srclient.Schema) error {

	orgBytes, err := json.Marshal(&org)
	if err != nil {
		return errors.Wrap(err, "json.Marshsal")
	}
	schemaIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(schemaIDBytes, uint32(schema.ID()))

	native, _, _ := schema.Codec().NativeFromTextual(orgBytes)
	valueBytes, err := schema.Codec().BinaryFromNative(nil, native)

	if err != nil {
		return errors.Wrap(err, "schema.Codec().BinaryFromNative")
	}

	// Combine row bytes with schema id to make a record
	var recordValue []byte
	recordValue = append(recordValue, byte(0))
	recordValue = append(recordValue, schemaIDBytes...)
	recordValue = append(recordValue, valueBytes...)
	key, _ := uuid.NewUUID()

	return o.orgProducer.PublishCreate(ctx, kafka.Message{
		Key:   []byte(key.String()),
		Value: recordValue,
		Time:  time.Now().UTC(),
	})
}

func (o *organizationUC) PublishUpdate(ctx context.Context, org *models.Organization) error {
	orgBytes, err := json.Marshal(&org)

	if err != nil {
		return errors.Wrap(err, "json.Marshal")
	}

	return o.orgProducer.PublishUpdate(ctx, kafka.Message{
		Value: orgBytes,
		Time:  time.Now().UTC(),
	})
}
