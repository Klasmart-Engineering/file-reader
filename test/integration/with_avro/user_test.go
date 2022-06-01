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
	"github.com/icrowley/fake"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestAvroConsumeUserCsv(t *testing.T) {
	// set up env variables
	userAvroTopic := "userAvroTopic" + uuid.NewString()
	s3FileCreationTopic := "s3FileCreatedTopic" + uuid.NewString()
	closer := env.EnvSetter(map[string]string{
		"USER_AVRO_TOPIC":                  userAvroTopic,
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
	bucket := "user"
	s3key := "user" + uuid.NewString() + ".csv"
	operationType := "user"

	// Make test csv file
	numUsers := 5
	userGeneratorMap := map[string]func() string{
		"uuid":               util.UuidFieldGenerator(),
		"user_given_name":    util.HumanNameFieldGenerator(2, 10),
		"user_family_name":   util.HumanNameFieldGenerator(2, 10),
		"user_email":         fake.EmailAddress,
		"user_phone_number":  fake.Phone,
		"user_date_of_birth": util.DateGenerator(1950, 2022, "2006-01-02"),
		"user_gender":        util.GenderGenerator(),
	}
	file, users := util.MakeCsv(numUsers, userGeneratorMap)

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
		Topic:       userAvroTopic,
		StartOffset: kafka.FirstOffset,
	})
	for i := 0; i < numUsers; i++ {
		msg, err := r.ReadMessage(ctx)
		assert.Nil(t, err, "error reading message from topic")
		userOutput, err := avro.DeserializeUser(bytes.NewReader(msg.Value[5:]))
		assert.Nil(t, err, "error deserialising message to user")
		t.Log(userOutput)

		assert.Equal(t, trackingId, userOutput.Metadata.Tracking_id)

		userInput := users[i]
		assert.Equal(t, userInput["uuid"], userOutput.Payload.Uuid)
		assert.Equal(t, userInput["user_given_name"], userOutput.Payload.Given_name)
		assert.Equal(t, userInput["user_family_name"], userOutput.Payload.Family_name)
		assert.Equal(t, userInput["user_email"], userOutput.Payload.Email)
		assert.Equal(t, userInput["user_date_of_birth"], userOutput.Payload.Date_of_birth)
		assert.Equal(t, userInput["user_gender"], userOutput.Payload.Gender)
		assert.Equal(t, userInput["user_phone_number"], userOutput.Payload.Phone_number.String)
	}
	ctx.Done()
}