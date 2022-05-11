package src

import (
	"bytes"
	"context"
	"encoding/csv"
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

type ConsumeToIngestConfig struct {
	BrokerAddrs     []string
	AwsSession      *session.Session
	OutputDirectory string
	Logger          log.Logger
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

func ConsumeToIngest(ctx context.Context, kafkaReader *kafka.Reader, config ConsumeToIngestConfig, operationMap map[string]Operation) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Read file-create message off kafka topic
			msg, err := kafkaReader.ReadMessage(ctx)
			if err != nil {
				panic("could not read message " + err.Error())
			}
			fmt.Println("received: ", string(msg.Value))

			// Deserialize file create message
			r := bytes.NewReader(msg.Value[5:]) // need to make consumer which checks/pulls schema from registry
			s3FileCreated, err := avro.DeserializeS3FileCreated(r) // then pass actual schema into the _FromSchema version
			if err != nil {
				panic("could not deserialize message " + err.Error())
			}

			// For now have the s3 downloader write to disk
			file, err := os.Create(config.OutputDirectory + s3FileCreated.Payload.Key)
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
			// Reopen the same file for ingest (until thought of alternative)
			f, _ := os.Open(config.OutputDirectory + s3FileCreated.Payload.Key)
			defer f.Close()

			// Compose the ingestFile() with different reader depending on file type
			var reader Reader
			switch s3FileCreated.Payload.Content_type {
			default:
				reader = csv.NewReader(f)
			}

			// Map to operation based on operation type
			operation, exists := operationMap[s3FileCreated.Payload.Operation_type]
			kafkaWriter := kafka.Writer{
				Addr:   kafka.TCP(config.BrokerAddrs...),
				Topic:  operation.Topic,
				Logger: &config.Logger,
			}
			if !exists {
				panic("invalid operation_type on file create message ")
			}

			operation.IngestFile(ctx, reader, kafkaWriter)
		}
	}
}
