package test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	avro "file_reader/avro_gencode"
	src "file_reader/src"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/riferrei/srclient"
	"github.com/segmentio/kafka-go"
)

func CreateBucket(sess *session.Session, bucket string) {
	svc := s3.New(sess)
	_, err := svc.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		panic(fmt.Sprintf("Unable to create bucket %q, %v", bucket, err))
	}
	fmt.Printf("Waiting for bucket %q to be created...\n", bucket)
	err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		panic(fmt.Sprintf("Error occurred while waiting for bucket to be created, %v", bucket))
	}
	fmt.Printf("Bucket %q successfully created\n", bucket)
}

func TestConsume(t *testing.T) {
	//t.Skip("f")
	schemaRegistryClient := &src.SchemaRegistry{
		C: srclient.CreateSchemaRegistryClient("http://localhost:8081"),
	}
	s3FileCreationTopic := "s3filecreation"
	schemaBody := avro.S3FileCreated.Schema(avro.NewS3FileCreated())
	s3FileCreationSchemaId := schemaRegistryClient.GetSchemaIdBytes(schemaBody, s3FileCreationTopic)
	kafkakey := "fakekey"
	brokerAddrs := []string{"localhost:9092"}

	awsRegion := "eu-west-1"

	bucket := "testbucket"
	fileDir := "./data/good/"
	filename := "organization.csv"
	s3key := filename

	ctx := context.Background()

	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "localstack",
		Config: aws.Config{
			Credentials:      credentials.NewStaticCredentials("fakekey", "fakesecretkey", ""),
			Region:           aws.String(awsRegion),
			Endpoint:         aws.String("http://localhost:4566"),
			S3ForcePathStyle: aws.Bool(true),
		},
	})
	if err != nil {
		fmt.Printf("Failed to initialize new session: %v", err)
		return
	}

	//CreateBucket(sess, bucket)

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

	//Combine row bytes with schema id to make a record
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

	// Start the Kafka consumer -> ingest file process
	operationMap := src.CreateOperationMap(schemaRegistryClient)
	var config = src.ConsumeToIngestConfig{
		BrokerAddrs:     brokerAddrs,
		AwsSession:      sess,
		OutputDirectory: "./data/downloaded/",
		Logger:          *log.New(os.Stdout, "kafka writer: ", 0),
	}
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokerAddrs,
		Topic:   s3FileCreationTopic,
	})

	src.ConsumeToIngest(ctx, r, config, operationMap)

}
