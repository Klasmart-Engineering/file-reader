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

func testProtoConsumeClassDetailsCsv(t *testing.T, numClassDetails int, classDetailsGeneratorMap map[string]func() string) {
	// set up env variables
	classDetailsProtoTopic := "classDetailsProtoTopic" + uuid.NewString()
	s3FileCreationTopic := "s3FileCreatedTopic" + uuid.NewString()
	closer := env.EnvSetter(map[string]string{
		"CLASS_DETAILS_PROTO_TOPIC":        classDetailsProtoTopic,
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
	bucket := "class-details"
	operationType := "class-details"
	s3key := "class-details" + uuid.NewString() + ".csv"

	// Make test csv file
	file, classDetails := util.MakeCsv(numClassDetails, classDetailsGeneratorMap)

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
		Brokers:     []string{"localhost:9092"},
		GroupID:     "consumer-group-" + uuid.NewString(),
		Topic:       classDetailsProtoTopic,
		StartOffset: kafka.FirstOffset,
	})

	serde := protobuf.NewProtoSerDe()
	classDetailsOutput := &onboarding.ClassDetails{}

	for i := 0; i < numClassDetails; i++ {
		msg, err := r.ReadMessage(ctx)
		assert.Nil(t, err, "error reading message from topic")

		_, err = serde.Deserialize(msg.Value, classDetailsOutput)

		assert.Nil(t, err, "error deserializing message from topic")

		assert.Equal(t, trackingUuid, classDetailsOutput.Metadata.TrackingUuid)

		classDetailsInput := classDetails[i]
		assert.Equal(t, classDetailsInput["class_uuid"], classDetailsOutput.Payload.ClassUuid)
		assert.Equal(t, classDetailsInput["school_uuid"], util.DerefString(classDetailsOutput.Payload.SchoolUuid))

		program_uuids := strings.Split(classDetailsInput["program_uuids"], ";")
		assert.Equal(t, program_uuids, classDetailsOutput.Payload.ProgramUuids)
		subject_uuids := strings.Split(classDetailsInput["subject_uuids"], ";")
		assert.Equal(t, subject_uuids, classDetailsOutput.Payload.SubjectUuids)
		grade_uuids := strings.Split(classDetailsInput["grade_uuids"], ";")
		assert.Equal(t, grade_uuids, classDetailsOutput.Payload.GradeUuids)
		age_range_uuids := strings.Split(classDetailsInput["age_range_uuids"], ";")
		assert.Equal(t, age_range_uuids, classDetailsOutput.Payload.AgeRangeUuids)

		assert.Equal(t, classDetailsInput["academic_term_uuid"], util.DerefString(classDetailsOutput.Payload.AcademicTermUuid))
	}
	ctx.Done()
}

func TestAvroConsumeClassDetailsCsvScenarios(t *testing.T) {
	type TestCases struct {
		description              string
		numClassDetails          int
		classDetailsGeneratorMap map[string]func() string
	}

	for _, scenario := range []TestCases{
		{
			description:     "should ingest class details when all optional fields are supplied",
			numClassDetails: 5,
			classDetailsGeneratorMap: map[string]func() string{
				"class_uuid":         util.UuidFieldGenerator(),
				"school_uuid":        util.UuidFieldGenerator(),
				"program_uuids":      util.RepeatedFieldGenerator(util.UuidFieldGenerator(), 5, 10),
				"subject_uuids":      util.RepeatedFieldGenerator(util.UuidFieldGenerator(), 5, 10),
				"grade_uuids":        util.RepeatedFieldGenerator(util.UuidFieldGenerator(), 5, 10),
				"age_range_uuids":    util.RepeatedFieldGenerator(util.UuidFieldGenerator(), 5, 10),
				"academic_term_uuid": util.UuidFieldGenerator(),
			},
		},
		{
			description:     "should ingest class details when all optional fields are null",
			numClassDetails: 5,
			classDetailsGeneratorMap: map[string]func() string{
				"class_uuid":         util.UuidFieldGenerator(),
				"school_uuid":        util.EmptyFieldGenerator(),
				"program_uuids":      util.EmptyFieldGenerator(),
				"subject_uuids":      util.EmptyFieldGenerator(),
				"grade_uuids":        util.EmptyFieldGenerator(),
				"age_range_uuids":    util.EmptyFieldGenerator(),
				"academic_term_uuid": util.EmptyFieldGenerator(),
			},
		},
	} {
		t.Run(scenario.description, func(t *testing.T) {
			testProtoConsumeClassDetailsCsv(t, scenario.numClassDetails, scenario.classDetailsGeneratorMap)
		})
	}
}
