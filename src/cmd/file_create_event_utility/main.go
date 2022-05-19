// Utility so that testers can put file creation events onto file_creation topic

package main

import (
	"bytes"
	"context"
	avro "file_reader/avro_gencode"
	"file_reader/src"
	"fmt"
	"log"
	"os"

	"github.com/riferrei/srclient"
	"github.com/segmentio/kafka-go"
)

var s3key = ""              // This should match to the file uploaded to S3
var trackingId = ""         // This will appear in the organization messages produced by the code, it should be the same coming out as going in
var bucket = "organization" // This can be anything as long as it is created beforehand

var awsRegion = os.Getenv("AWS_DEFAULT_REGION")
var content_type = "text/csv" // Current implementation attempts to ingest csv as a default case, and there are not other file types implemented, so this won't change behavior

func main() {
	ctx := context.Background()

	// Get schema id from registry
	schemaRegistryClient := &src.SchemaRegistry{
		C: srclient.CreateSchemaRegistryClient(os.Getenv("SCHEMA_CLIENT_ENDPOINT")),
	}
	s3FileCreationTopic := "S3FileCreatedUpdated"
	schemaBody := avro.S3FileCreated.Schema(avro.NewS3FileCreated())
	s3FileCreationSchemaId := schemaRegistryClient.GetSchemaIdBytes(schemaBody, s3FileCreationTopic)

	// Encode file_created message using schema
	s3FileCreatedCodec := avro.S3FileCreated{
		Payload: avro.S3FileCreatedPayload{
			Key:            s3key,
			Aws_region:     awsRegion,
			Bucket_name:    bucket,
			Content_length: 0, // Content length isn't yet implemented
			Content_type:   content_type,
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

	// Put the message on the file_created topic
	brokerAddrs := []string{os.Getenv("BROKERS")}
	w := kafka.Writer{
		Addr:                   kafka.TCP(brokerAddrs...),
		Topic:                  s3FileCreationTopic,
		AllowAutoTopicCreation: true,
		Logger:                 log.New(os.Stdout, "kafka writer: ", 0),
	}
	err := w.WriteMessages(
		ctx,
		kafka.Message{
			Value: recordValue,
		},
	)
	if err != nil {
		fmt.Println(err)
	}
}
