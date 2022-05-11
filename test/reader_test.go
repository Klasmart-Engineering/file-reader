package test

import (
	"bytes"
	"context"
	avro "file_reader/avro_gencode"
	"file_reader/src"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/riferrei/srclient"
	"github.com/segmentio/kafka-go"
)

func TestConsume(t *testing.T) {
	// This will be refactored into an integration test
	t.Skip("f")
	schemaRegistryClient := &src.SchemaRegistry{
		C: srclient.CreateSchemaRegistryClient(os.Getenv("SCHEMA_CLIENT_ENDPOINT")),
	}
	s3FileCreationTopic := os.Getenv("S3_FILE_CREATED_UPDATED_TOPIC")
	schemaBody := avro.S3FileCreated.Schema(avro.NewS3FileCreated())
	s3FileCreationSchemaId := schemaRegistryClient.GetSchemaIdBytes(schemaBody, s3FileCreationTopic)
	kafkakey := ""
	brokerAddrs := []string{os.Getenv("KAFKA_BROKER")}

	awsRegion := os.Getenv("AWS_DEFAULT_REGION")

	bucket := "organization"
	fileDir := "./data/good/"
	filename := "organization.csv"
	s3key := filename

	ctx := context.Background()

	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "localstack",
		Config: aws.Config{
			Credentials: credentials.NewStaticCredentials(
				os.Getenv("AWS_ACCESS_KEY_ID"),
				os.Getenv("AWS_SECRET_ACCESS_KEY"),
				"",
			),
			Region:           aws.String(os.Getenv("AWS_DEFAULT_REGION")),
			Endpoint:         aws.String(os.Getenv("AWS_ENDPOINT")),
			S3ForcePathStyle: aws.Bool(true),
		},
	})
	if err != nil {
		fmt.Printf("Failed to initialize new session: %v", err)
		return
	}

	// Upload file to s3
	file, err := os.Open(fileDir + filename)
	if err != nil {
		panic(fmt.Sprintf("Unable to open file %q, %v", err, err))
	}
	defer file.Close()
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(s3key),
		Body:   file,
	})
	if err != nil {
		panic(fmt.Sprintf("Unable to upload %q to %q, %v", filename, bucket, err))
	}

	// Put file create message on topic
	fi, err := file.Stat()
	if err != nil {
		// Print the error and exit.
		panic(fmt.Sprintf("Unable to get file information %q, %v", filename, err))
	}
	s3FileCreatedCodec := avro.S3FileCreated{
		Payload: avro.S3FileCreatedPayload{
			Key:            s3key,
			Aws_region:     awsRegion,
			Bucket_name:    bucket,
			Content_length: fi.Size(),
			Content_type:   "text/csv",
			Operation_type: "organization",
		},
		Metadata: avro.S3FileCreatedMetadata{},
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
		Addr:   kafka.TCP(brokerAddrs...),
		Topic:  s3FileCreationTopic,
		Logger: log.New(os.Stdout, "kafka writer: ", 0),
	}
	err = w.WriteMessages(
		ctx,
		kafka.Message{
			Key:   []byte(kafkakey),
			Value: recordValue,
		},
	)
	if err != nil {
		panic("could not write message " + err.Error())
	}
}
