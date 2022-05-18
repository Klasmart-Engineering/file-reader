package integration_test

import (
	"bytes"
	"context"
	"encoding/binary"
	avro "file_reader/avro_gencode"
	"file_reader/src"
	"file_reader/src/instrument"
	"file_reader/src/pkg/proto"
	protoSrc "file_reader/src/third_party/protobuf/srclient"
	"file_reader/test/env"
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
	// set env
	closer := env.EnvSetter(map[string]string{
		"BROKERS":                  "localhost:9092",
		"PROTO_SCHEMA_DIRECTORY":   "protos/onboarding",
		"SCHEMA_CLIENT_ENDPOINT":   "http://localhost:8081",
		"ORGANIZATION_PROTO_TOPIC": uuid.NewString(),
		"DOWNLOAD_DIRECTORY":       "./" + uuid.NewString(),
	})

	defer t.Cleanup(closer) // In Go 1.14+

	// proto schema registry
	protoSRC := proto.GetNewSchemaRegistry(
		protoSrc.NewClient(protoSrc.WithURL(instrument.MustGetEnv("SCHEMA_CLIENT_ENDPOINT"))),
		context.Background(),
	)

	// avros schema registry
	schemaRegistryClient := &src.SchemaRegistry{
		C: srclient.CreateSchemaRegistryClient(instrument.MustGetEnv("SCHEMA_CLIENT_ENDPOINT")),
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

	// Read file-create message off kafka topic
	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokerAddrs,
		GroupID:     os.Getenv("ORGANIZATION_GROUP_ID"),
		StartOffset: kafka.LastOffset,
		Topic:       os.Getenv("S3_FILE_CREATED_UPDATED_TOPIC"),
	})

	msg, err := kafkaReader.ReadMessage(ctx)
	if err != nil {
		t.Errorf("could not read message ", err.Error())
	}
	t.Log("received message: ", string(msg.Value))

	// Deserialize file create message
	schemaIdBytes := msg.Value[1:5]
	schemaId := int(binary.BigEndian.Uint32(schemaIdBytes))
	schema, err := schemaRegistryClient.GetSchema(schemaId)
	if err != nil {
		t.Errorf("could not retrieve schema with id ", schemaId, err.Error())
	}
	r := bytes.NewReader(msg.Value[5:])
	s3FileCreated, err := avro.DeserializeS3FileCreatedFromSchema(r, schema)
	if err != nil {
		t.Errorf("could not deserialize message ", err.Error())
	}

	// For now have the s3 downloader write to disk
	file, err := os.Create(config.OutputDirectory + s3FileCreated.Payload.Key)
	if err != nil {
		logger.Error(ctx, err)
		continue
	}
	defer file.Close()

	for i := 0; i < numOrgs; i++ {
		msg, err := r.ReadMessage(ctx)
		assert.Nil(t, err, "error reading message from topic")
		orgOutput, err := avro.DeserializeOrganization(bytes.NewReader(msg.Value[5:]))
		assert.Nil(t, err, "error deserialising message to org")
		t.Log(orgOutput)

		assert.Equal(t, trackingId, orgOutput.Metadata.Tracking_id)

		orgInput := orgs[i]
		assert.Equal(t, orgInput[0], orgOutput.Payload.Guid)
		assert.Equal(t, orgInput[1], orgOutput.Payload.Organization_name)
	}
}
