package integration_test

import (
	"context"
	"strings"
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

func testProtoConsumeSchoolCsv(t *testing.T, numSchools int, schoolGeneratorMap map[string]func() string) {
	// set up env variables
	schoolProtoTopic := "schoolProtoTopic" + uuid.NewString()
	s3FileCreationTopic := "s3FileCreatedTopic" + uuid.NewString()
	closer := env.EnvSetter(map[string]string{
		"SCHOOL_PROTO_TOPIC":               schoolProtoTopic,
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
	bucket := "school"
	operationType := "school"
	s3key := "school" + uuid.NewString() + ".csv"

	// Make test csv file
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
		Topic:       schoolProtoTopic,
		StartOffset: kafka.FirstOffset,
	})

	serde := protobuf.NewProtoSerDe()
	schoolOutput := &onboarding.School{}

	for i := 0; i < numSchools; i++ {
		msg, err := r.ReadMessage(ctx)
		assert.Nil(t, err, "error reading message from topic")

		_, err = serde.Deserialize(msg.Value, schoolOutput)

		assert.Nil(t, err, "error deserializing message from topic")

		assert.Equal(t, trackingId, schoolOutput.Metadata.TrackingId)

		schoolInput := schools[i]
		assert.Equal(t, schoolInput["uuid"], util.DerefString(schoolOutput.Payload.Uuid))
		assert.Equal(t, schoolInput["school_name"], schoolOutput.Payload.Name)
		assert.Equal(t, schoolInput["organization_id"], schoolOutput.Payload.OrganizationId)
		program_ids := strings.Split(schoolInput["program_ids"], ";")
		assert.Equal(t, program_ids, util.DerefArrayString(&schoolOutput.Payload.ProgramIds))
	}
	ctx.Done()
}

func TestAvroConsumeSchoolCsvScenarios(t *testing.T) {
	type TestCases struct {
		description        string
		numSchools         int
		schoolGeneratorMap map[string]func() string
	}

	for _, scenario := range []TestCases{
		{
			description: "should ingest schools when all optional fields are supplied",
			numSchools:  5,
			schoolGeneratorMap: map[string]func() string{
				"uuid":            util.UuidFieldGenerator(),
				"organization_id": util.UuidFieldGenerator(),
				"school_name":     util.NameFieldGenerator("school", 5),
				"program_ids":     util.RepeatedFieldGenerator(util.UuidFieldGenerator(), 5, 10),
			},
		},
		{
			description: "should ingest schools when all optional fields are null",
			numSchools:  5,
			schoolGeneratorMap: map[string]func() string{
				"uuid":            util.EmptyFieldGenerator(),
				"organization_id": util.UuidFieldGenerator(),
				"school_name":     util.NameFieldGenerator("school", 5),
				"program_ids":     util.EmptyFieldGenerator(),
			},
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			testProtoConsumeSchoolCsv(t, scenario.numSchools, scenario.schoolGeneratorMap)
		})
	}
}
