package integration_test

import (
	"bytes"
	"context"
	"encoding/binary"
	"log"
	"os"
	"testing"

	avro "github.com/KL-Engineering/file-reader/api/avro/avro_gencode"
	"github.com/KL-Engineering/file-reader/internal/core"
	zapLogger "github.com/KL-Engineering/file-reader/internal/log"
	"github.com/KL-Engineering/file-reader/test/env"
	util "github.com/KL-Engineering/file-reader/test/integration"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/riferrei/srclient"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestConsumeS3CsvOrganization(t *testing.T) {
	// set up env variables
	organizationAvroTopic := "orgAvroTopic" + uuid.NewString()
	s3FileCreationTopic := "s3FileCreatedTopic" + uuid.NewString()
	closer := env.EnvSetter(map[string]string{
		"ORGANIZATION_AVRO_TOPIC":          organizationAvroTopic,
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

	schemaRegistryClient := &core.SchemaRegistry{
		C: srclient.CreateSchemaRegistryClient("http://localhost:8081"),
	}
	schemaBody := avro.S3FileCreated.Schema(avro.NewS3FileCreated())
	s3FileCreationSchemaId := schemaRegistryClient.GetSchemaId(schemaBody, s3FileCreationTopic)
	schemaIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(schemaIDBytes, uint32(s3FileCreationSchemaId))
	kafkakey := ""
	brokerAddrs := []string{"localhost:9092"}

	awsRegion := "eu-west-1"

	bucket := "organization"
	s3key := "organization" + uuid.NewString() + ".csv"

	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Credentials: credentials.NewStaticCredentials(
				"test",
				"test",
				"",
			),
			Region:           aws.String("eu-west-1"),
			Endpoint:         aws.String("http://localhost:4566"),
			S3ForcePathStyle: aws.Bool(true),
		},
	})
	assert.Nil(t, err, "error creating aws session")

	// Upload file to s3

	numOrgs := 5
	orgGeneratorMap := map[string]func() string{
		"uuid":              util.UuidFieldGenerator(),
		"owner_user_id":     util.UuidFieldGenerator(),
		"id_list":           util.RepeatedFieldGenerator(util.UuidFieldGenerator(), 0, 5),
		"foo":               util.UuidFieldGenerator(),
		"bar":               util.UuidFieldGenerator(),
		"organization_name": util.NameFieldGenerator("org", numOrgs),
	}

	file, orgs := util.MakeCsv(numOrgs, orgGeneratorMap)
	//file, orgs := util.MakeOrgsCsv(numOrgs)
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(s3key),
		Body:   file,
	})
	assert.Nil(t, err, "error putting s3 object to bucket")

	// Put file create message on topic
	trackingId := uuid.NewString()
	s3FileCreatedCodec := avro.S3FileCreated{
		Payload: avro.S3FileCreatedPayload{
			Key:            s3key,
			Aws_region:     awsRegion,
			Bucket_name:    bucket,
			Content_length: 0, // Content length isn't yet implemented
			Content_type:   "text/csv",
			Operation_type: "organization",
		},
		Metadata: avro.S3FileCreatedMetadata{Tracking_id: trackingId},
	}
	var buf bytes.Buffer
	s3FileCreatedCodec.Serialize(&buf)
	valueBytes := buf.Bytes()

	// Combine row bytes with schema id to make a record
	var recordValue []byte
	recordValue = append(recordValue, byte(0))
	recordValue = append(recordValue, schemaIDBytes...)
	recordValue = append(recordValue, valueBytes...)

	w := kafka.Writer{
		Addr:                   kafka.TCP(brokerAddrs...),
		Topic:                  s3FileCreationTopic,
		AllowAutoTopicCreation: true,
		Logger:                 log.New(os.Stdout, "kafka writer: ", 0),
	}
	err = w.WriteMessages(
		ctx,
		kafka.Message{
			Key:   []byte(kafkakey),
			Value: recordValue,
		},
	)
	assert.Nil(t, err, "error writing message to topic")

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{"localhost:9092"},
		GroupID:     "consumer-group-" + uuid.NewString(),
		Topic:       organizationAvroTopic,
		StartOffset: kafka.FirstOffset,
		Dialer: &kafka.Dialer{
			Timeout: kafka.DefaultDialer.Timeout,
		},
	})
	for i := 0; i < numOrgs; i++ {
		msg, err := r.ReadMessage(ctx)
		assert.Nil(t, err, "error reading message from topic")
		orgOutput, err := avro.DeserializeOrganization(bytes.NewReader(msg.Value[5:]))
		assert.Nil(t, err, "error deserialising message to org")
		t.Log(orgOutput)

		assert.Equal(t, trackingId, orgOutput.Metadata.Tracking_id)

		orgInput := orgs[i]
		assert.Equal(t, orgInput["uuid"], orgOutput.Payload.Uuid)
		assert.Equal(t, orgInput["organization_name"], orgOutput.Payload.Name)
		assert.Equal(t, orgInput["owner_user_id"], orgOutput.Payload.Owner_user_id)
	}
}
