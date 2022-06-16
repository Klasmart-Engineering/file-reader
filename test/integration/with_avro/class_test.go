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

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func testAvrosConsumeClassCsv(t *testing.T, numClasses int, classGeneratorMap map[string]func() string) {
	// set up env variables
	classAvroTopic := "classAvroTopic" + uuid.NewString()
	s3FileCreationTopic := "s3FileCreatedTopic" + uuid.NewString()
	closer := env.EnvSetter(map[string]string{
		"CLASS_AVRO_TOPIC":                 classAvroTopic,
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
	s3key := "class" + uuid.NewString() + ".csv"
	operationType := "class"

	// Make test csv file
	file, classes := util.MakeCsv(numClasses, classGeneratorMap)

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

	// Assert that the real file got ingested
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokerAddrs,
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

		assert.Equal(t, trackingUuid, classOutput.Metadata.Tracking_uuid)

		classInput := classes[i]
		assert.Equal(t, classInput["uuid"], util.DerefAvroNullString(classOutput.Payload.Uuid))
		assert.Equal(t, classInput["name"], classOutput.Payload.Name)
		assert.Equal(t, classInput["organization_uuid"], classOutput.Payload.Organization_uuid)
	}
	ctx.Done()

}
func TestAvroConsumeClassCsvScenarios(t *testing.T) {

	// test cases
	testCases := []struct {
		name              string
		numClasses        int
		classGeneratorMap map[string]func() string
	}{
		{
			name:       "Should ingest classes when all optional fields are supplied",
			numClasses: 5,
			classGeneratorMap: map[string]func() string{
				"uuid":              util.UuidFieldGenerator(),
				"organization_uuid": util.UuidFieldGenerator(),
				"name":              util.NameFieldGenerator("class", 5),
			},
		},
		{
			name:       "Should ingest classes when all optional fields are null",
			numClasses: 5,
			classGeneratorMap: map[string]func() string{
				"uuid":              util.EmptyFieldGenerator(),
				"organization_uuid": util.UuidFieldGenerator(),
				"name":              util.NameFieldGenerator("class", 5),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testAvrosConsumeClassCsv(t, tc.numClasses, tc.classGeneratorMap)
		})
	}

}
func TestAvroConsumeInvalidAndValidClassCsv(t *testing.T) {
	// set up env variables
	classAvroTopic := "classAvroTopic" + uuid.NewString()
	s3FileCreationTopic := "s3FileCreatedTopic" + uuid.NewString()
	closer := env.EnvSetter(map[string]string{
		"CLASS_AVRO_TOPIC":                 classAvroTopic,
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

	// Then try to consume a real class file
	s3key2 := "organization" + uuid.NewString() + ".csv"
	numClasses := 5
	classGeneratorMap := map[string]func() string{
		"uuid":              util.UuidFieldGenerator(),
		"organization_uuid": util.UuidFieldGenerator(),
		"fake_uuid":         util.UuidFieldGenerator(),
		"name":              util.NameFieldGenerator("class", numClasses),
	}
	file, classes := util.MakeCsv(numClasses, classGeneratorMap)
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
		Topic:       classAvroTopic,
		StartOffset: kafka.FirstOffset,
	})
	for i := 0; i < numClasses; i++ {
		msg, err := r.ReadMessage(ctx)
		assert.Nil(t, err, "error reading message from topic")
		classOutput, err := avro.DeserializeClass(bytes.NewReader(msg.Value[5:]))
		assert.Nil(t, err, "error deserialising message to class")
		t.Log(classOutput)

		assert.Equal(t, trackingUuid2, classOutput.Metadata.Tracking_uuid)

		classInput := classes[i]
		assert.Equal(t, classInput["uuid"], classOutput.Payload.Uuid.String)
		assert.Equal(t, classInput["name"], classOutput.Payload.Name)
		assert.Equal(t, classInput["organization_uuid"], classOutput.Payload.Organization_uuid)
	}
	ctx.Done()
}
