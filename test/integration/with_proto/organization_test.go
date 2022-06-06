package integration_test

import (
	"context"

	avro "github.com/KL-Engineering/file-reader/api/avro/avro_gencode"
	"github.com/KL-Engineering/file-reader/api/proto/proto_gencode/onboarding"
	"github.com/KL-Engineering/file-reader/internal/core"
	zapLogger "github.com/KL-Engineering/file-reader/internal/log"
	util "github.com/KL-Engineering/file-reader/test/integration"

	"github.com/KL-Engineering/file-reader/pkg/third_party/protobuf"
	"github.com/KL-Engineering/file-reader/test/env"

	"testing"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestProtoConsumeOrganizationCsv(t *testing.T) {
	// set up env variables
	organizationProtoTopic := "orgProtoTopic" + uuid.NewString()
	s3FileCreationTopic := "s3FileCreatedTopic" + uuid.NewString()
	closer := env.EnvSetter(map[string]string{
		"ORGANIZATION_PROTO_TOPIC":         organizationProtoTopic,
		"S3_FILE_CREATED_UPDATED_GROUP_ID": "s3FileCreatedGroupId" + uuid.NewString(),
		"S3_FILE_CREATED_UPDATED_TOPIC":    s3FileCreationTopic,
		"SCHEMA_TYPE":                      "PROTO",
	})

	defer t.Cleanup(closer)

	// Start consumer
	l, _ := zap.NewDevelopment()
	logger := zapLogger.Wrap(l)
	ctx := context.Background()
	core.StartFileCreateConsumer(ctx, logger)

	brokerAddrs := []string{"localhost:9092"}
	awsRegion := "eu-west-1"
	bucket := "organization"
	operationType := "organization"
	s3key := "organization" + uuid.NewString() + ".csv"

	// Make test csv file
	numOrgs := 5
	orgGeneratorMap := map[string]func() string{
		"uuid":              util.UuidFieldGenerator(),
		"owner_user_id":     util.UuidFieldGenerator(),
		"id_list":           util.RepeatedFieldGenerator(util.UuidFieldGenerator(), 0, 5),
		"foo":               util.UuidFieldGenerator(),
		"bar":               util.UuidFieldGenerator(),
		"organization_name": util.NameFieldGenerator("org", numOrgs),
	}

	file, orgs := util.MakeCsv(numOrgs, orgGeneratorMap)

	// Upload csv to S3
	err := util.UploadFileToS3(bucket, s3key, awsRegion, file)
	assert.Nil(t, err, "error uploading file to s3")

	// Put file create message on topic
	trackingId := uuid.NewString()
	s3FileCreated := avro.S3FileCreated{
		Payload: avro.S3FileCreatedPayload{
			Key:            s3key,
			Aws_region:     awsRegion,
			Bucket_name:    bucket,
			Content_length: 0, // Content length isn't yet implemented
			Content_type:   "text/csv",
			Operation_type: operationType,
		},
		Metadata: avro.S3FileCreatedMetadata{Tracking_id: trackingId},
	}
	err = util.ProduceFileCreateMessage(
		ctx,
		s3FileCreationTopic,
		brokerAddrs,
		s3FileCreated,
	)
	assert.Nil(t, err, "error producing file create message to topic")

	// Consume from output topic and make assertions
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{"localhost:9092"},
		GroupID:     "consumer-group-" + uuid.NewString(),
		Topic:       organizationProtoTopic,
		StartOffset: kafka.FirstOffset,
	})

	serde := protobuf.NewProtoSerDe()
	orgOutput := &onboarding.Organization{}

	for i := 0; i < numOrgs; i++ {
		msg, err := r.ReadMessage(ctx)
		assert.Nil(t, err, "error reading message from topic")

		_, err = serde.Deserialize(msg.Value, orgOutput)

		assert.Nil(t, err, "error deserializing message from topic")

		assert.Equal(t, trackingId, orgOutput.Metadata.TrackingId)

		orgInput := orgs[i]
		assert.Equal(t, orgInput["uuid"], orgOutput.Payload.Uuid)
		assert.Equal(t, orgInput["organization_name"], orgOutput.Payload.Name)
		assert.Equal(t, orgInput["owner_user_id"], orgOutput.Payload.OwnerUserId)
	}
	ctx.Done()
}
