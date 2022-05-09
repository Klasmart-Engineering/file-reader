package src

import (
	"bytes"
	"context"
	avro "file_reader/avro_gencode"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/segmentio/kafka-go"
)

// <- also generally need to think about how all this will work given:
// there are multiple operations and it needs to find the right one
// there could be other file types than csv

type ConsumeS3Config struct {
	Topic       string
	BrokerAddrs []string
	AwsSession  *session.Session
	Logger      log.Logger
}

func CreateOperationMap(schemaRegistryClient *SchemaRegistry) map[string]Operation {
	return map[string]Operation{
		"organization": {
			Topic:         OrganizationTopic,
			Key:           "",
			SchemaIDBytes: GetOrganizationSchemaIdBytes(schemaRegistryClient),
			RowToSchema:   RowToOrganization,
		},
	}
}

func ConsumeToIngest(ctx context.Context, config ConsumeS3Config, operationMap map[string]Operation) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: config.BrokerAddrs,
		Topic:   config.Topic,
	})
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Read file create message off kafka topic
			msg, err := r.ReadMessage(ctx) // test that ctx works for kill
			if err != nil {
				panic("could not read message " + err.Error())
			}
			fmt.Println("received: ", string(msg.Value))

			// Deserialize file create message
			r := bytes.NewReader(msg.Value[5:])
			s3FileCreated, err := avro.DeserializeS3FileCreated(r)
			if err != nil {
				panic("could not deserialize message " + err.Error())
			}

			// For now have s3 download write to disk
			file, err := os.Create("./data/downloaded/" + s3FileCreated.Payload.Key)
			if err != nil {
				fmt.Println(err)
			}
			defer file.Close()

			downloader := s3manager.NewDownloader(config.AwsSession)
			numBytes, err := downloader.Download(file,
				&s3.GetObjectInput{
					Bucket: aws.String(s3FileCreated.Payload.Bucket_name),
					Key:    aws.String(s3FileCreated.Payload.Key),
				})
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println("Downloaded", s3FileCreated.Payload.Key, numBytes, "bytes")
			file.Close()

			// Ingest downloaded file
			f, _ := os.Open("./data/downloaded/" + s3FileCreated.Payload.Key)
			defer f.Close()

			var ingestConfig = IngestConfig{
				BrokerAddrs: config.BrokerAddrs,
				Reader:      f,
				Context:     context.Background(),
				Logger:      config.Logger,
			}
			operation, exists := operationMap[s3FileCreated.Payload.Operation_type]
			if !exists {
				panic("invalid operation_type on file create message ")
			}
			operation.IngestFile(ingestConfig, s3FileCreated.Payload.Content_type)
		}
	}
}
