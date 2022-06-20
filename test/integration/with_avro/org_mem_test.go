package integration_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	avro "github.com/KL-Engineering/file-reader/api/avro/avro_gencode"
	"github.com/KL-Engineering/file-reader/internal/core"
	zapLogger "github.com/KL-Engineering/file-reader/internal/log"
	"github.com/KL-Engineering/file-reader/test/env"
	util "github.com/KL-Engineering/file-reader/test/integration"
	"go.uber.org/zap"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

func TestAvroConsumeOrgMemCsv(t *testing.T) {
	// set up env variables
	orgMemAvroTopic := "orgMemAvroTopic" + uuid.NewString()
	s3FileCreationTopic := "s3FileCreatedTopic" + uuid.NewString()
	closer := env.EnvSetter(map[string]string{
		"ORGANIZATION_MEMBERSHIP_AVRO_TOPIC": orgMemAvroTopic,
		"S3_FILE_CREATED_UPDATED_GROUP_ID":   "s3FileCreatedGroupId" + uuid.NewString(),
		"S3_FILE_CREATED_UPDATED_TOPIC":      s3FileCreationTopic,
		"SCHEMA_TYPE":                        "AVRO",
	})

	defer t.Cleanup(closer)

	ctx := context.Background()
	// Start consumer
	l, _ := zap.NewDevelopment()
	logger := zapLogger.Wrap(l)
	core.StartFileCreateConsumer(ctx, logger)

	brokerAddrs := []string{"localhost:9092"}
	awsRegion := "eu-west-1"
	bucket := "organization-membership"
	s3key := "organization-membership" + uuid.NewString() + ".csv"
	operationType := "organization_membership"

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

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokerAddrs,
		GroupID:     "consumer-group-" + uuid.NewString(),
		Topic:       orgMemAvroTopic,
		StartOffset: kafka.FirstOffset,
	})
	for i := 0; i < numOrgMems; i++ {
		msg, err := r.ReadMessage(ctx)
		assert.Nil(t, err, "error reading message from topic")
		orgMemOutput, err := avro.DeserializeOrganizationMembership(bytes.NewReader(msg.Value[5:]))
		assert.Nil(t, err, "error deserialising message to org")
		t.Log(orgMemOutput)

		assert.Equal(t, trackingUuid, orgMemOutput.Metadata.Tracking_uuid)

		orgMemInput := orgMems[i]
		assert.Equal(t, orgMemInput["organization_uuid"], orgMemOutput.Payload.Organization_uuid)
		assert.Equal(t, orgMemInput["user_uuid"], orgMemOutput.Payload.User_uuid)
		orgRoleUuids := strings.Split(orgMemInput["organization_role_uuids"], ";")
		assert.Equal(t, orgRoleUuids, orgMemOutput.Payload.Organization_role_uuids)
	}
	ctx.Done()
}

func TestAvroConsumeInvalidAndValidOrgMemCsv(t *testing.T) {
	t.Skip()
	// set up env variables
	orgMemAvroTopic := "orgMemAvroTopic" + uuid.NewString()
	s3FileCreationTopic := "s3FileCreatedTopic" + uuid.NewString()
	closer := env.EnvSetter(map[string]string{
		"ORGANIZATION_MEMBERSHIP_AVRO_TOPIC": orgMemAvroTopic,
		"S3_FILE_CREATED_UPDATED_GROUP_ID":   "s3FileCreatedGroupId" + uuid.NewString(),
		"S3_FILE_CREATED_UPDATED_TOPIC":      s3FileCreationTopic,
		"SCHEMA_TYPE":                        "AVRO",
	})

	defer t.Cleanup(closer)

	ctx := context.Background()
	// Start consumer
	l, _ := zap.NewDevelopment()
	logger := zapLogger.Wrap(l)
	core.StartFileCreateConsumer(ctx, logger)

	brokerAddrs := []string{"localhost:9092"}
	awsRegion := "eu-west-1"
	bucket := "organization-membership"
	operationType := "organization_membership"

	// First try to consume an empty file
	s3key1 := "bad_organization_membership" + uuid.NewString() + ".csv"
	emptyFile := util.MakeEmptyFile()
	err := util.UploadFileToS3(bucket, s3key1, awsRegion, emptyFile)
	assert.Nil(t, err, "error uploading file to s3")
	trackingUuid1 := uuid.NewString()
	s3FileCreated1 := avro.S3FileCreatedUpdated{
		Payload: avro.S3FileCreatedUpdatedPayload{
			Key:            s3key1,
			Aws_region:     awsRegion,
			Bucket_name:    bucket,
			Content_length: 0, // Content length isn't yet implemented
			Content_type:   "text/csv",
			Operation_type: operationType,
		},
		Metadata: avro.S3FileCreatedUpdatedMetadata{Tracking_uuid: trackingUuid1},
	}
	err = util.ProduceFileCreateMessage(
		ctx,
		s3FileCreationTopic,
		brokerAddrs,
		s3FileCreated1,
	)
	assert.Nil(t, err, "error producing file create message to topic")

	// Then try to consume a real organization file
	s3key2 := "organization-membership" + uuid.NewString() + ".csv"
	numOrgMems := 5
	orgMemGeneratorMap := map[string]func() string{
		"organization_uuid":       util.UuidFieldGenerator(),
		"user_uuid":               util.UuidFieldGenerator(),
		"organization_role_uuids": util.RepeatedFieldGenerator(util.UuidFieldGenerator(), 0, numOrgMems),
	}
	file, orgMems := util.MakeCsv(numOrgMems, orgMemGeneratorMap)
	err = util.UploadFileToS3(bucket, s3key2, awsRegion, file)
	assert.Nil(t, err, "error uploading file to s3")
	trackingUuid2 := uuid.NewString()
	s3FileCreated2 := avro.S3FileCreatedUpdated{
		Payload: avro.S3FileCreatedUpdatedPayload{
			Key:            s3key2,
			Aws_region:     awsRegion,
			Bucket_name:    bucket,
			Content_length: 0, // Content length isn't yet implemented
			Content_type:   "text/csv",
			Operation_type: operationType,
		},
		Metadata: avro.S3FileCreatedUpdatedMetadata{Tracking_uuid: trackingUuid2},
	}
	err = util.ProduceFileCreateMessage(
		ctx,
		s3FileCreationTopic,
		brokerAddrs,
		s3FileCreated2,
	)
	assert.Nil(t, err, "error producing file create message to topic")

	// Assert that the real file got ingested
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokerAddrs,
		GroupID:     "consumer-group-" + uuid.NewString(),
		Topic:       orgMemAvroTopic,
		StartOffset: kafka.FirstOffset,
	})
	for i := 0; i < numOrgMems; i++ {
		msg, err := r.ReadMessage(ctx)
		assert.Nil(t, err, "error reading message from topic")
		orgMemOutput, err := avro.DeserializeOrganizationMembership(bytes.NewReader(msg.Value[5:]))
		assert.Nil(t, err, "error deserialising message to org membership")
		t.Log(orgMemOutput)

		assert.Equal(t, trackingUuid2, orgMemOutput.Metadata.Tracking_uuid)

		orgMemInput := orgMems[i]
		assert.Equal(t, orgMemInput["organization_uuid"], orgMemOutput.Payload.Organization_uuid)
		assert.Equal(t, orgMemInput["user_uuid"], orgMemOutput.Payload.User_uuid)

		orgRoleUuids := strings.Split(orgMemInput["organization_role_uuids"], ";")
		assert.Equal(t, orgRoleUuids, orgMemOutput.Payload.Organization_role_uuids)
	}
	ctx.Done()

}
