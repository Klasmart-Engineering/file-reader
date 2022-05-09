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
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/riferrei/srclient"
	"github.com/segmentio/kafka-go"
)

// Skip for now
func TestReadOrgCsv(t *testing.T) {
	t.Skip("Skip for now")

	schemaRegistryClient := &src.SchemaRegistry{
		C: srclient.CreateSchemaRegistryClient("http://localhost:8081"),
	}

	// Fake CSV data
	orgId1 := uuid.NewString()
	orgId2 := uuid.NewString()
	orgName1 := "org_name1"
	orgName2 := "org_name2"
	csv := fmt.Sprintf("%s,%s\n%s,%s", orgId1, orgName1, orgId2, orgName2)
	reader := bytes.NewReader([]byte(csv))

	brokerAddrs := []string{"localhost:9092"}

	var config = src.IngestConfig{
		BrokerAddrs: brokerAddrs,
		Reader:      reader,
		Context:     context.Background(),
		Logger:      *log.New(os.Stdout, "kafka writer: ", 0),
	}
	Organization := src.Operation{
		Topic:         src.OrganizationTopic,
		Key:           "",
		SchemaIDBytes: src.GetOrganizationSchemaIdBytes(schemaRegistryClient),
		RowToSchema:   src.RowToOrganization,
	}
	Organization.IngestFile(config, "text/csv")

	// Until we have a fresh topic for testing,
	// Not sure yet how to do assertions as consumer will have to read everything
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokerAddrs,
		Topic:   "organization",
	})
	ctx := context.Background()
	for i := 0; i < 20; i++ {
		msg, err := r.ReadMessage(ctx)
		if err == nil {
			val, e := avro.DeserializeOrganization(bytes.NewReader(msg.Value[5:]))
			if e == nil {
				t.Logf("Here is the message %s\n", val)
			} else {
				t.Logf("Error deserializing: %e", e)
			}
		} else {
			t.Logf("Error consuming the message: %v (%v)\n", err, msg)
			break
		}
	}
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func TestConsume(t *testing.T) {
	//t.Skip("f")

	os.Setenv("METADATA_ORIGIN_APPLICATION", "")
	os.Setenv("METADATA_REGION", "")
	schemaRegistryClient := &src.SchemaRegistry{
		C: srclient.CreateSchemaRegistryClient("http://localhost:8081"),
	}
	s3FileCreationTopic := "s3filecreation"
	schemaBody := avro.S3FileCreated.Schema(avro.NewS3FileCreated())
	s3FileCreationSchemaId := schemaRegistryClient.GetSchemaIdBytes(schemaBody, s3FileCreationTopic)
	awsRegion := "eu-west-1"
	kafkakey := "fakekey"

	ctx := context.Background()
	// Might be okay to keep sess code here in refactor
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
	// REMOVE ALL BELOW
	bucket := "testbucket"
	filename := "organization.csv"
	s3key := filename
	// Create bucket <- include as a test utility function?
	// svc := s3.New(sess)
	// _, err = svc.CreateBucket(&s3.CreateBucketInput{
	// 	Bucket: aws.String(bucket),
	// })
	// if err != nil {
	// 	exitErrorf("Unable to create bucket %q, %v", bucket, err)
	// }
	// fmt.Printf("Waiting for bucket %q to be created...\n", bucket)
	// err = svc.WaitUntilBucketExists(&s3.HeadBucketInput{
	// 	Bucket: aws.String(bucket),
	// })
	// if err != nil {
	// 	exitErrorf("Error occurred while waiting for bucket to be created, %v", bucket)
	// }
	// fmt.Printf("Bucket %q successfully created\n", bucket)

	// Upload file to s3 <- include as a test utlity function?
	file, err := os.Open("./data/good/" + filename)
	if err != nil {
		exitErrorf("Unable to open file %q, %v", err)
	}
	defer file.Close()
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(s3key),
		Body:   file,
	})
	if err != nil {
		// Print the error and exit.
		exitErrorf("Unable to upload %q to %q, %v", filename, bucket, err)
	}
	brokerAddrs := []string{"localhost:9092"}
	fi, err := file.Stat()
	if err != nil {
		// Print the error and exit.
		exitErrorf("Unable to get file information %q, %v", filename, err)
	}
	var config = src.ConsumeS3Config{
		Topic:       s3FileCreationTopic,
		BrokerAddrs: brokerAddrs,
		AwsSession:  sess,
		Logger:      *log.New(os.Stdout, "kafka writer: ", 0),
	}
	// Put message onto s3file topic
	var buf bytes.Buffer
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
	operationMap := src.CreateOperationMap(schemaRegistryClient)
	src.ConsumeToIngest(ctx, config, operationMap) // <- don't know why doesn't work as goroutine
	// <- still not sure how we will test it
}
