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

func testAvroConsumeClassDetailsCsv(t *testing.T, numClassDetails int, classDetailsGeneratorMap map[string]func() string) {
	// set up env variables
	classDetailsAvroTopic := "classDetailsAvroTopic" + uuid.NewString()
	s3FileCreationTopic := "s3FileCreatedTopic" + uuid.NewString()
	closer := env.EnvSetter(map[string]string{
		"CLASS_DETAILS_AVRO_TOPIC":         classDetailsAvroTopic,
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
	bucket := "class-details"
	s3key := "class-details" + uuid.NewString() + ".csv"
	operationType := "class_details"

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
		Topic:       classDetailsAvroTopic,
		StartOffset: kafka.FirstOffset,
	})
	for i := 0; i < numClassDetails; i++ {
		msg, err := r.ReadMessage(ctx)
		assert.Nil(t, err, "error reading message from topic")
		classDetailsOutput, err := avro.DeserializeClassDetails(bytes.NewReader(msg.Value[5:]))
		assert.Nil(t, err, "error deserialising message to class details")
		t.Log(classDetailsOutput)

		assert.Equal(t, trackingUuid, classDetailsOutput.Metadata.Tracking_uuid)

		classDetailsInput := classDetails[i]
		assert.Equal(t, classDetailsInput["class_uuid"], classDetailsOutput.Payload.Class_uuid)
		assert.Equal(t, classDetailsInput["school_uuid"], util.DerefAvroNullString(classDetailsOutput.Payload.School_uuid))
		program_ids := strings.Split(classDetailsInput["program_uuids"], ";")
		assert.Equal(t, program_ids, util.DerefAvroNullArrayString(classDetailsOutput.Payload.Program_uuids))
		subject_uuids := strings.Split(classDetailsInput["subject_uuids"], ";")
		assert.Equal(t, subject_uuids, util.DerefAvroNullArrayString(classDetailsOutput.Payload.Subject_uuids))
		grade_uuids := strings.Split(classDetailsInput["grade_uuids"], ";")
		assert.Equal(t, grade_uuids, util.DerefAvroNullArrayString(classDetailsOutput.Payload.Grade_uuids))
		age_range_uuids := strings.Split(classDetailsInput["age_range_uuids"], ";")
		assert.Equal(t, age_range_uuids, util.DerefAvroNullArrayString(classDetailsOutput.Payload.Age_range_uuids))
		assert.Equal(t, classDetailsInput["academic_term_uuid"], util.DerefAvroNullString(classDetailsOutput.Payload.Academic_term_uuid))
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
			testAvroConsumeClassDetailsCsv(t, scenario.numClassDetails, scenario.classDetailsGeneratorMap)
		})
	}
}
