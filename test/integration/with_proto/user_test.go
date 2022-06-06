package integration_test

import (
	"context"

	avro "github.com/KL-Engineering/file-reader/api/avro/avro_gencode"
	"github.com/KL-Engineering/file-reader/api/proto/proto_gencode/onboarding"
	"github.com/KL-Engineering/file-reader/internal/core"
	zapLogger "github.com/KL-Engineering/file-reader/internal/log"
	util "github.com/KL-Engineering/file-reader/test/integration"
	"github.com/icrowley/fake"

	"github.com/KL-Engineering/file-reader/pkg/third_party/protobuf"
	"github.com/KL-Engineering/file-reader/test/env"

	"testing"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestProtoConsumeUserCsv(t *testing.T) {
	// set up env variables
	userProtoTopic := "userProtoTopic" + uuid.NewString()
	s3FileCreationTopic := "s3FileCreatedTopic" + uuid.NewString()
	closer := env.EnvSetter(map[string]string{
		"USER_PROTO_TOPIC":                 userProtoTopic,
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
	bucket := "user"
	operationType := "user"
	s3key := "user" + uuid.NewString() + ".csv"

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

	// Consume from output topic and make assertions
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{"localhost:9092"},
		GroupID:     "consumer-group-" + uuid.NewString(),
		Topic:       userProtoTopic,
		StartOffset: kafka.FirstOffset,
	})

	serde := protobuf.NewProtoSerDe()
	userOutput := &onboarding.User{}

	for i := 0; i < numUsers; i++ {
		msg, err := r.ReadMessage(ctx)
		assert.Nil(t, err, "error reading message from topic")

		_, err = serde.Deserialize(msg.Value, userOutput)

		assert.Nil(t, err, "error deserializing message from topic")

		assert.Equal(t, trackingId, userOutput.Metadata.TrackingId)

		userInput := users[i]
		assert.Equal(t, userInput["uuid"], userOutput.Payload.Uuid)
		assert.Equal(t, userInput["user_given_name"], userOutput.Payload.GivenName)
		assert.Equal(t, userInput["user_family_name"], userOutput.Payload.FamilyName)
		assert.Equal(t, userInput["user_email"], util.DerefString(userOutput.Payload.Email))
		assert.Equal(t, userInput["user_phone_number"], util.DerefString(userOutput.Payload.PhoneNumber))
		assert.Equal(t, userInput["user_date_of_birth"], util.DerefString(userOutput.Payload.DateOfBirth))
		assert.Equal(t, userInput["user_gender"], userOutput.Payload.Gender)
	}
	ctx.Done()
}
