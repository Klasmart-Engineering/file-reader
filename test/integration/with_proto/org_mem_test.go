package integration_test

import (
	"context"
	"strings"

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

func TestProtoConsumeOrgMemCsv(t *testing.T) {
	// set up env variables
	orgMemProtoTopic := "orgMemProtoTopic" + uuid.NewString()
	s3FileCreationTopic := "s3FileCreatedTopic" + uuid.NewString()
	closer := env.EnvSetter(map[string]string{
		"ORGANIZATION_MEMBERSHIP_PROTO_TOPIC": orgMemProtoTopic,
		"S3_FILE_CREATED_UPDATED_GROUP_ID":    "s3FileCreatedGroupId" + uuid.NewString(),
		"S3_FILE_CREATED_UPDATED_TOPIC":       s3FileCreationTopic,
		"SCHEMA_TYPE":                         "PROTO",
	})

	defer t.Cleanup(closer)

	// Start consumer
	l, _ := zap.NewDevelopment()
	logger := zapLogger.Wrap(l)
	ctx := context.Background()
	core.StartFileCreateConsumer(ctx, logger)

	brokerAddrs := []string{"localhost:9092"}
	awsRegion := "eu-west-1"
	bucket := "organization.membership"
	operationType := "organization_membership"
	s3key := "organization.membership" + uuid.NewString() + ".csv"

	// Make test csv file
	numOrgMems := 5
	orgMemGeneratorMap := map[string]func() string{
		"organization_uuid":       util.UuidFieldGenerator(),
		"user_uuid":               util.UuidFieldGenerator(),
		"organization_role_uuids": util.RepeatedFieldGenerator(util.UuidFieldGenerator(), 0, numOrgMems),
	}

	file, orgMems := util.MakeCsv(numOrgMems, orgMemGeneratorMap)

	// Upload csv to S3
	err := util.UploadFileToS3(bucket, s3key, awsRegion, file)
	assert.Nil(t, err, "error uploading file to s3")

	// Put file create message on topic
	trackingUuid := uuid.NewString()
	s3FileCreated := avro.S3FileCreatedUpdated{
		Payload: avro.S3FileCreatedUpdatedPayload{
			Key:            s3key,
			Aws_region:     awsRegion,
			Bucket_name:    bucket,
			Content_length: 0, // Content length isn't yet implemented
			Content_type:   "text/csv",
			Operation_type: operationType,
		},
		Metadata: avro.S3FileCreatedUpdatedMetadata{Tracking_uuid: trackingUuid},
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
		Topic:       orgMemProtoTopic,
		StartOffset: kafka.FirstOffset,
	})

	serde := protobuf.NewProtoSerDe()
	orgMemOutput := &onboarding.OrganizationMembership{}

	for i := 0; i < numOrgMems; i++ {
		msg, err := r.ReadMessage(ctx)
		assert.Nil(t, err, "error reading message from topic")

		_, err = serde.Deserialize(msg.Value, orgMemOutput)

		assert.Nil(t, err, "error deserializing message from topic")

		assert.Equal(t, trackingUuid, orgMemOutput.Metadata.TrackingUuid)

		orgMemInput := orgMems[i]
		assert.Equal(t, orgMemInput["organization_uuid"], orgMemOutput.Payload.OrganizationUuid)
		assert.Equal(t, orgMemInput["user_uuid"], orgMemOutput.Payload.UserUuid)

		orgRoleUuids := strings.Split(orgMemInput["organization_role_uuids"], ";")
		assert.Equal(t, orgRoleUuids, orgMemOutput.Payload.OrganizationRoleUuids)
	}
	ctx.Done()
}
