package integration_test

import (
	"bytes"
	"context"
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

func TestAvroConsumeSchoolMemCsv(t *testing.T) {
	// set up env variables
	schoolMemAvroTopic := "schoolMemAvroTopic" + uuid.NewString()
	s3FileCreationTopic := "s3FileCreatedTopic" + uuid.NewString()
	closer := env.EnvSetter(map[string]string{
		"SCHOOL_MEMBERSHIP_AVRO_TOPIC":     schoolMemAvroTopic,
		"S3_FILE_CREATED_UPDATED_GROUP_ID": "s3FileCreatedGroupId" + uuid.NewString(),
		"S3_FILE_CREATED_UPDATED_TOPIC":    s3FileCreationTopic,
		"SCHEMA_TYPE":                      "AVRO",
	})

	defer t.Cleanup(closer)

	ctx := context.Background()
	// Start consumer
	l, _ := zap.NewDevelopment()
	logger := zapLogger.Wrap(l)
	core.StartFileCreateConsumer(ctx, logger)

	brokerAddrs := []string{"localhost:9092"}
	awsRegion := "eu-west-1"
	bucket := "school.membership"
	s3key := "school.membership" + uuid.NewString() + ".csv"
	operationType := "school_membership"

	// Make test csv file
	numSchoolMems := 5
	schoolMemGeneratorMap := map[string]func() string{
		"school_uuid": util.UuidFieldGenerator(),
		"user_uuid":   util.UuidFieldGenerator(),
	}

	file, schoolMems := util.MakeCsv(numSchoolMems, schoolMemGeneratorMap)

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
		Topic:       schoolMemAvroTopic,
		StartOffset: kafka.FirstOffset,
	})
	for i := 0; i < numSchoolMems; i++ {
		msg, err := r.ReadMessage(ctx)
		assert.Nil(t, err, "error reading message from topic")
		schoolMemOutput, err := avro.DeserializeSchoolMembership(bytes.NewReader(msg.Value[5:]))
		assert.Nil(t, err, "error deserialising message to school membership")
		t.Log(schoolMemOutput)

		assert.Equal(t, trackingUuid, schoolMemOutput.Metadata.Tracking_uuid)

		schoolMemInput := schoolMems[i]
		assert.Equal(t, schoolMemInput["school_uuid"], schoolMemOutput.Payload.School_uuid)
		assert.Equal(t, schoolMemInput["user_uuid"], schoolMemOutput.Payload.User_uuid)
	}
	ctx.Done()
}

func TestAvroConsumeInvalidAndValidSchoolMemCsv(t *testing.T) {
	// set up env variables
	schoolMemAvroTopic := "schoolMemAvroTopic" + uuid.NewString()
	s3FileCreationTopic := "s3FileCreatedTopic" + uuid.NewString()
	closer := env.EnvSetter(map[string]string{
		"SCHOOL_MEMBERSHIP_AVRO_TOPIC":     schoolMemAvroTopic,
		"S3_FILE_CREATED_UPDATED_GROUP_ID": "s3FileCreatedGroupId" + uuid.NewString(),
		"S3_FILE_CREATED_UPDATED_TOPIC":    s3FileCreationTopic,
		"SCHEMA_TYPE":                      "AVRO",
	})

	defer t.Cleanup(closer)

	ctx := context.Background()
	// Start consumer
	l, _ := zap.NewDevelopment()
	logger := zapLogger.Wrap(l)
	core.StartFileCreateConsumer(ctx, logger)

	brokerAddrs := []string{"localhost:9092"}
	awsRegion := "eu-west-1"
	bucket := "school.membership"
	operationType := "school_membership"

	// First try to consume an empty file
	s3key1 := "bad_school_membership" + uuid.NewString() + ".csv"
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
	s3key2 := "organization.membership" + uuid.NewString() + ".csv"
	numSchoolMems := 5
	schoolMemGeneratorMap := map[string]func() string{
		"school_uuid": util.UuidFieldGenerator(),
		"user_uuid":   util.UuidFieldGenerator(),
	}
	file, schoolMems := util.MakeCsv(numSchoolMems, schoolMemGeneratorMap)
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
		Topic:       schoolMemAvroTopic,
		StartOffset: kafka.FirstOffset,
	})
	for i := 0; i < numSchoolMems; i++ {
		msg, err := r.ReadMessage(ctx)
		assert.Nil(t, err, "error reading message from topic")
		schoolMemOutput, err := avro.DeserializeSchoolMembership(bytes.NewReader(msg.Value[5:]))
		assert.Nil(t, err, "error deserialising message to school membership")
		t.Log(schoolMemOutput)

		assert.Equal(t, trackingUuid2, schoolMemOutput.Metadata.Tracking_uuid)

		schoolMemInput := schoolMems[i]
		assert.Equal(t, schoolMemInput["school_uuid"], schoolMemOutput.Payload.School_uuid)
		assert.Equal(t, schoolMemInput["user_uuid"], schoolMemOutput.Payload.User_uuid)
	}
	ctx.Done()

}