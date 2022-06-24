package integration_test

import (
	"context"
	"testing"

	avro "github.com/KL-Engineering/file-reader/api/avro/avro_gencode"
	"github.com/KL-Engineering/file-reader/api/proto/proto_gencode/onboarding"
	"github.com/KL-Engineering/file-reader/internal/core"
	zapLogger "github.com/KL-Engineering/file-reader/internal/log"
	"github.com/KL-Engineering/file-reader/pkg/third_party/protobuf"
	"github.com/KL-Engineering/file-reader/test/env"
	util "github.com/KL-Engineering/file-reader/test/integration"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func testProtoConsumeClassRosterCsv(t *testing.T, numClassRosters int, classRosterGeneratorMap map[string]func() string) {
	// set up env variables
	classRosterProtoTopic := "classRosterProtoTopic" + uuid.NewString()
	s3FileCreationTopic := "s3FileCreatedTopic" + uuid.NewString()
	closer := env.EnvSetter(map[string]string{
		"CLASS_ROSTER_PROTO_TOPIC":         classRosterProtoTopic,
		"S3_FILE_CREATED_UPDATED_GROUP_ID": "s3FileCreatedGroupId" + uuid.NewString(),
		"S3_FILE_CREATED_UPDATED_TOPIC":    s3FileCreationTopic,
		"SCHEMA_TYPE":                      "PROTO",
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
	s3key := "classroster" + uuid.NewString() + ".csv"

	// Make test csv file

	file, classRosters := util.MakeCsv(numClassRosters, classRosterGeneratorMap)

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
			Content_length: file.Size(),
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
		Topic:       classRosterProtoTopic,
		StartOffset: kafka.FirstOffset,
	})

	serde := protobuf.NewProtoSerDe()
	classRosterOutput := &onboarding.ClassRoster{}

	for i := 0; i < numClassRosters; i++ {
		msg, err := r.ReadMessage(ctx)
		assert.Nil(t, err, "error reading message from topic")

		_, err = serde.Deserialize(msg.Value, classRosterOutput)

		assert.Nil(t, err, "error deserializing message from topic")

		assert.Equal(t, trackingUuid, classRosterOutput.Metadata.TrackingUuid)

		classRosterInput := classRosters[i]
		assert.Equal(t, classRosterInput["class_uuid"], classRosterOutput.Payload.ClassUuid)
		assert.Equal(t, classRosterInput["user_uuid"], classRosterOutput.Payload.UserUuid)
		assert.Equal(t, classRosterInput["participating_as"], classRosterOutput.Payload.ParticipatingAs)

	}
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
			testProtoConsumeClassRosterCsv(t, tc.numClassRosters, tc.classRosterGeneratorMap)
		})
	}

}
func TestProtoConsumeInvalidAndValidClassRosterCsv(t *testing.T) {
	// set up env variables
	classRosterProtoTopic := "classRosterProtoTopic" + uuid.NewString()
	s3FileCreationTopic := "s3FileCreatedTopic" + uuid.NewString()
	closer := env.EnvSetter(map[string]string{
		"CLASS_ROSTER_PROTO_TOPIC":         classRosterProtoTopic,
		"S3_FILE_CREATED_UPDATED_GROUP_ID": "s3FileCreatedGroupId" + uuid.NewString(),
		"S3_FILE_CREATED_UPDATED_TOPIC":    s3FileCreationTopic,
		"SCHEMA_TYPE":                      "PROTO",
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
			Content_length: emptyFile.Size(),
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
			Content_length: file.Size(),
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
		Topic:       classRosterProtoTopic,
		StartOffset: kafka.FirstOffset,
	})

	serde := protobuf.NewProtoSerDe()
	classRosterOutput := &onboarding.ClassRoster{}

	for i := 0; i < numClassRosters; i++ {
		msg, err := r.ReadMessage(ctx)
		assert.Nil(t, err, "error reading message from topic")

		_, err = serde.Deserialize(msg.Value, classRosterOutput)

		assert.Nil(t, err, "error deserializing message from topic")

		assert.Equal(t, trackingUuid2, classRosterOutput.Metadata.TrackingUuid)

		classRosterInput := classRosters[i]
		assert.Equal(t, classRosterInput["class_uuid"], classRosterOutput.Payload.ClassUuid)
		assert.Equal(t, classRosterInput["user_uuid"], classRosterOutput.Payload.UserUuid)
		assert.Equal(t, classRosterInput["participating_as"], classRosterOutput.Payload.ParticipatingAs)

	}
}
