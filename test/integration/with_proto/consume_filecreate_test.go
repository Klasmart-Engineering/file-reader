package integration_test

import (
	"bytes"
	"context"
	avro "file_reader/avro_gencode"
	"file_reader/src"
	"file_reader/src/instrument"
	zapLogger "file_reader/src/log"
	"file_reader/src/pkg/validation"
	"file_reader/src/protos/onboarding"
	"file_reader/src/third_party/protobuf"
	util "file_reader/test/integration"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

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

func MakeOrgsCsv(numOrgs int) (csv *strings.Reader, orgs [][]string) {
	organizations := [][]string{}
	for i := 0; i < numOrgs; i++ {
		// rows are `uuid,org{i}`
		organizations = append(organizations, []string{uuid.NewString(), "org" + strconv.Itoa(i)})
	}
	lines := []string{}
	for _, org := range organizations {
		s := strings.Join(org, ",")
		lines = append(lines, s)
	}
	file := strings.NewReader(strings.Join(lines, "\n"))
	return file, organizations
}

func TestConsumeS3CsvOrganization(t *testing.T) {

	schemaRegistryClient := &src.SchemaRegistry{
		C: srclient.CreateSchemaRegistryClient("http://localhost:8081"),
	}
	s3FileCreationTopic := "S3FileCreatedUpdated"
	schemaBody := avro.S3FileCreated.Schema(avro.NewS3FileCreated())
	s3FileCreationSchemaId := schemaRegistryClient.GetSchemaIdBytes(schemaBody, s3FileCreationTopic)
	kafkakey := ""
	brokerAddrs := []string{"localhost:9092"}

	awsRegion := "eu-west-1"

	bucket := "organization"
	s3key := "organization.csv"

	ctx := context.Background()

	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "localstack",
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
	file, orgs := MakeOrgsCsv(numOrgs)
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
	recordValue = append(recordValue, s3FileCreationSchemaId...)
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

	l, _ := zap.NewDevelopment()
	logger := zapLogger.Wrap(l)

	// start consumer
	go util.StartFileCreateConsumer(ctx, logger)

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{"localhost:9092"},
		GroupID:     "consumer-group-" + uuid.NewString(),
		Topic:       instrument.MustGetEnv("ORGANIZATION_PROTO_TOPIC"),
		StartOffset: kafka.LastOffset,
	})

	serde := protobuf.NewProtoSerDe()
	orgOutput := &onboarding.Organization{}

	for i := 0; i < numOrgs; i++ {
		msg, err := r.ReadMessage(ctx)
		assert.Nil(t, err, "error reading message from topic")

		_, err = serde.Deserialize(msg.Value, orgOutput)

		assert.Nil(t, err, "error deserializing message from topic")

		t.Log(orgOutput)

		validateTrackingId := validation.ValidateTrackingId{Uuid: orgOutput.Metadata.TrackingId.Value}

		err = validation.UUIDValidate(validateTrackingId)

		assert.Nil(t, err, "UUID is invalid")

		orgInput := orgs[i]
		assert.Equal(t, orgInput[0], orgOutput.Payload.Uuid.Value)
		assert.Equal(t, orgInput[1], orgOutput.Payload.Name.Value)
	}
}
