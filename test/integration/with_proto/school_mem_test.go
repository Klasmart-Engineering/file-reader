package integration_test

import (
	"context"

	avro "github.com/KL-Engineering/file-reader/api/avro/avro_gencode"
	"github.com/KL-Engineering/file-reader/api/proto/proto_gencode/onboarding"
	"github.com/KL-Engineering/file-reader/internal/core"
	zapLogger "github.com/KL-Engineering/file-reader/internal/log"
	util "github.com/KL-Engineering/file-reader/test/integration"

	"github.com/KL-Engineering/file-reader/pkg/third_party/protobuf"
	"github.com/KL-Engineering/file-reader/test/env"

	"testing"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestProtoConsumeSchoolMemCsv(t *testing.T) {
	// set up env variables
	schoolMemProtoTopic := "schoolMemProtoTopic" + uuid.NewString()
	s3FileCreationTopic := "s3FileCreatedTopic" + uuid.NewString()
	closer := env.EnvSetter(map[string]string{
		"SCHOOL_MEMBERSHIP_PROTO_TOPIC":    schoolMemProtoTopic,
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
	bucket := "school.membership"
	operationType := "school_membership"
	s3key := "school.membership" + uuid.NewString() + ".csv"

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

	// Consume from output topic and make assertions
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{"localhost:9092"},
		GroupID:     "consumer-group-" + uuid.NewString(),
		Topic:       schoolMemProtoTopic,
		StartOffset: kafka.FirstOffset,
	})

	serde := protobuf.NewProtoSerDe()
	schoolMemOutput := &onboarding.SchoolMembership{}

	for i := 0; i < numSchoolMems; i++ {
		msg, err := r.ReadMessage(ctx)
		assert.Nil(t, err, "error reading message from topic")

		_, err = serde.Deserialize(msg.Value, schoolMemOutput)

		assert.Nil(t, err, "error deserializing message from topic")

		assert.Equal(t, trackingUuid, schoolMemOutput.Metadata.TrackingUuid)

		schoolMemInput := schoolMems[i]
		assert.Equal(t, schoolMemInput["school_uuid"], schoolMemOutput.Payload.SchoolUuid)
		assert.Equal(t, schoolMemInput["user_uuid"], schoolMemOutput.Payload.UserUuid)
	}
	ctx.Done()
}
