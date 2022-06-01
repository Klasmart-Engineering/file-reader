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
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestAvroConsumeSchoolCsv(t *testing.T) {
	// set up env variables
	schoolAvroTopic := "schoolAvroTopic" + uuid.NewString()
	s3FileCreationTopic := "s3FileCreatedTopic" + uuid.NewString()
	closer := env.EnvSetter(map[string]string{
		"SCHOOL_AVRO_TOPIC":                schoolAvroTopic,
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
	bucket := "school"
	s3key := "school" + uuid.NewString() + ".csv"
	operationType := "school"

	// Make test csv file
	numSchools := 5
	schoolGeneratorMap := map[string]func() string{
		"uuid":            util.UuidFieldGenerator(),
		"organization_id": util.UuidFieldGenerator(),
		"school_name":     util.NameFieldGenerator("school", numSchools),
		"program_ids":     util.RepeatedFieldGenerator(util.UuidFieldGenerator(), 5, 10),
		"fake_ids":        util.RepeatedFieldGenerator(util.UuidFieldGenerator(), 1, 5),
	}
	file, schools := util.MakeCsv(numSchools, schoolGeneratorMap)

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

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{"localhost:9092"},
		GroupID:     "consumer-group-" + uuid.NewString(),
		Topic:       schoolAvroTopic,
		StartOffset: kafka.FirstOffset,
	})
	for i := 0; i < numSchools; i++ {
		msg, err := r.ReadMessage(ctx)
		assert.Nil(t, err, "error reading message from topic")
		schoolOutput, err := avro.DeserializeSchool(bytes.NewReader(msg.Value[5:]))
		assert.Nil(t, err, "error deserialising message to school")
		t.Log(schoolOutput)

		assert.Equal(t, trackingId, schoolOutput.Metadata.Tracking_id)

		schoolInput := schools[i]
		assert.Equal(t, schoolInput["uuid"], schoolOutput.Payload.Uuid)
		assert.Equal(t, schoolInput["school_name"], schoolOutput.Payload.Name)
		assert.Equal(t, schoolInput["organization_id"], schoolOutput.Payload.Organization_id)
		program_ids := strings.Split(schoolInput["program_ids"], ";")
		assert.Equal(t, program_ids, schoolOutput.Payload.Program_ids)
	}
	ctx.Done()
}

func TestAvroConsumeInvalidAndValidSchoolCsv(t *testing.T) {
	// set up env variables
	schoolAvroTopic := "schoolAvroTopic" + uuid.NewString()
	s3FileCreationTopic := "s3FileCreatedTopic" + uuid.NewString()
	closer := env.EnvSetter(map[string]string{
		"SCHOOL_AVRO_TOPIC":                schoolAvroTopic,
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
	bucket := "class"
	operationType := "class"

	// First try to consume an empty file
	s3key1 := "bad_class" + uuid.NewString() + ".csv"
	emptyFile := util.MakeEmptyFile()
	err := util.UploadFileToS3(bucket, s3key1, awsRegion, emptyFile)
	assert.Nil(t, err, "error uploading file to s3")
	trackingId1 := uuid.NewString()
	s3FileCreated1 := avro.S3FileCreated{
		Payload: avro.S3FileCreatedPayload{
			Key:            s3key1,
			Aws_region:     awsRegion,
			Bucket_name:    bucket,
			Content_length: 0, // Content length isn't yet implemented
			Content_type:   "text/csv",
			Operation_type: operationType,
		},
		Metadata: avro.S3FileCreatedMetadata{Tracking_id: trackingId1},
	}
	err = util.ProduceFileCreateMessage(
		ctx,
		s3FileCreationTopic,
		brokerAddrs,
		s3FileCreated1,
	)
	assert.Nil(t, err, "error producing file create message to topic")

	// Then try to consume a real organization file
	s3key2 := "organization" + uuid.NewString() + ".csv"
	numClasses := 5
	classGeneratorMap := map[string]func() string{
		"uuid":            util.UuidFieldGenerator(),
		"organization_id": util.UuidFieldGenerator(),
		"id_list":         util.RepeatedFieldGenerator(util.UuidFieldGenerator(), 0, 5),
		"foo":             util.UuidFieldGenerator(),
		"bar":             util.UuidFieldGenerator(),
		"class_name":      util.NameFieldGenerator("class", numClasses),
	}
	file, classes := util.MakeCsv(numClasses, classGeneratorMap)
	err = util.UploadFileToS3(bucket, s3key2, awsRegion, file)
	assert.Nil(t, err, "error uploading file to s3")
	trackingId2 := uuid.NewString()
	s3FileCreated2 := avro.S3FileCreated{
		Payload: avro.S3FileCreatedPayload{
			Key:            s3key2,
			Aws_region:     awsRegion,
			Bucket_name:    bucket,
			Content_length: 0, // Content length isn't yet implemented
			Content_type:   "text/csv",
			Operation_type: operationType,
		},
		Metadata: avro.S3FileCreatedMetadata{Tracking_id: trackingId2},
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
		Brokers:     []string{"localhost:9092"},
		GroupID:     "consumer-group-" + uuid.NewString(),
		Topic:       classAvroTopic,
		StartOffset: kafka.FirstOffset,
	})
	for i := 0; i < numClasses; i++ {
		msg, err := r.ReadMessage(ctx)
		assert.Nil(t, err, "error reading message from topic")
		classOutput, err := avro.DeserializeClass(bytes.NewReader(msg.Value[5:]))
		assert.Nil(t, err, "error deserialising message to class")
		t.Log(classOutput)

		assert.Equal(t, trackingId2, classOutput.Metadata.Tracking_id)

		classInput := classes[i]
		assert.Equal(t, classInput["uuid"], classOutput.Payload.Uuid)
		assert.Equal(t, classInput["class_name"], classOutput.Payload.Name)
		assert.Equal(t, classInput["organization_id"], classOutput.Payload.Organization_uuid)
	}
	ctx.Done()
}
