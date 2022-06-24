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

func testAvrosConsumeClassRosterCsv(t *testing.T, numClassRosters int, classRosterGeneratorMap map[string]func() string) {
	// set up env variables
	classRosterAvroTopic := "classRosterAvroTopic" + uuid.NewString()
	s3FileCreationTopic := "s3FileCreatedTopic" + uuid.NewString()
	closer := env.EnvSetter(map[string]string{
		"CLASS_ROSTER_AVRO_TOPIC":          classRosterAvroTopic,
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
	bucket := "class-roster"
	s3key := "class-roster" + uuid.NewString() + ".csv"
	operationType := "class_roster"

	// Make test csv file
	file, classRosters := util.MakeCsv(numClassRosters, classRosterGeneratorMap)

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
		Topic:       classRosterAvroTopic,
		StartOffset: kafka.FirstOffset,
	})
	for i := 0; i < numClassRosters; i++ {
		msg, err := r.ReadMessage(ctx)
		assert.Nil(t, err, "error reading message from topic")
		classRosterOutput, err := avro.DeserializeClassRoster(bytes.NewReader(msg.Value[5:]))
		assert.Nil(t, err, "error deserialising message to class roster")
		t.Log(classRosterOutput)

		assert.Equal(t, trackingUuid, classRosterOutput.Metadata.Tracking_uuid)

		classRosterInput := classRosters[i]
		assert.Equal(t, classRosterInput["class_uuid"], classRosterOutput.Payload.Class_uuid)
		assert.Equal(t, classRosterInput["user_uuid"], classRosterOutput.Payload.User_uuid)
		assert.Equal(t, classRosterInput["participating_as"], classRosterOutput.Payload.Participating_as)
	}
	ctx.Done()

}
func TestAvroConsumeClassRosterCsvScenarios(t *testing.T) {

	// test cases
	testCases := []struct {
		name                    string
		numClassRosters         int
		classRosterGeneratorMap map[string]func() string
	}{
		{
			name:            "Should ingest class rosters when all compulsory fields are supplied",
			numClassRosters: 5,
			classRosterGeneratorMap: map[string]func() string{
				"class_uuid":       util.UuidFieldGenerator(),
				"user_uuid":        util.UuidFieldGenerator(),
				"participating_as": util.HumanNameFieldGenerator(2, 10),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testAvrosConsumeClassRosterCsv(t, tc.numClassRosters, tc.classRosterGeneratorMap)
		})
	}

}
func TestAvroConsumeInvalidAndValidClassRosterCsv(t *testing.T) {
	// set up env variables
	classRosterAvroTopic := "classRosterAvroTopic" + uuid.NewString()
	s3FileCreationTopic := "s3FileCreatedTopic" + uuid.NewString()
	closer := env.EnvSetter(map[string]string{
		"CLASS_ROSTER_AVRO_TOPIC":          classRosterAvroTopic,
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
	bucket := "class-roster"
	operationType := "class_roster"

	// First try to consume an empty file
	s3key1 := "bad_classroster" + uuid.NewString() + ".csv"
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
	s3key2 := "classroster" + uuid.NewString() + ".csv"
	numClassRosters := 5
	classRosterGeneratorMap := map[string]func() string{
		"class_uuid":       util.UuidFieldGenerator(),
		"user_uuid":        util.UuidFieldGenerator(),
		"fake_uuid":        util.UuidFieldGenerator(),
		"participating_as": util.HumanNameFieldGenerator(2, 10),
	}
	file, classRosters := util.MakeCsv(numClassRosters, classRosterGeneratorMap)
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
		Topic:       classRosterAvroTopic,
		StartOffset: kafka.FirstOffset,
	})
	for i := 0; i < numClassRosters; i++ {
		msg, err := r.ReadMessage(ctx)
		assert.Nil(t, err, "error reading message from topic")
		classRosterOutput, err := avro.DeserializeClassRoster(bytes.NewReader(msg.Value[5:]))
		assert.Nil(t, err, "error deserialising message to class roster")
		t.Log(classRosterOutput)

		assert.Equal(t, trackingUuid2, classRosterOutput.Metadata.Tracking_uuid)

		classRosterInput := classRosters[i]
		assert.Equal(t, classRosterInput["class_uuid"], classRosterOutput.Payload.Class_uuid)
		assert.Equal(t, classRosterInput["user_uuid"], classRosterOutput.Payload.User_uuid)
		assert.Equal(t, classRosterInput["participating_as"], classRosterOutput.Payload.Participating_as)
	}
	ctx.Done()
}
