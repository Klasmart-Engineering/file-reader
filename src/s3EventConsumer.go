package src

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/csv"
	avro "file_reader/avro_gencode"
	"file_reader/src/instrument"
	zaplogger "file_reader/src/log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/segmentio/kafka-go"
)

type ConsumeToIngestConfig struct {
	OutputBrokerAddrs []string
	AwsSession        *session.Session
	SchemaRegistry    *SchemaRegistry
	OperationMap      map[string]Operation
	OutputDirectory   string
	Logger            *zaplogger.ZapLogger
}

func CreateOperationMap(schemaRegistryClient *SchemaRegistry) map[string]Operation {
	// creates a map of key to Operation struct
	return map[string]Operation{
		"organization": {
			Topic:         OrganizationTopic,
			Key:           "",
			SchemaIDBytes: GetOrganizationSchemaIdBytes(schemaRegistryClient),
			RowToSchema:   RowToOrganization,
		},
	}
}

func ConsumeToIngest(ctx context.Context, kafkaReader *kafka.Reader, config ConsumeToIngestConfig) {
	logger := config.Logger
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Read file-create message off kafka topic
			msg, err := kafkaReader.ReadMessage(ctx)
			if err != nil {
				logger.Error(ctx, "could not read message ", err.Error())
				continue
			}
			logger.Info(ctx, "received message: ", string(msg.Value))

			// Deserialize file create message
			schemaIdBytes := msg.Value[1:5]
			schemaId := int(binary.BigEndian.Uint32(schemaIdBytes))
			schema, err := config.SchemaRegistry.GetSchema(schemaId)
			if err != nil {
				logger.Error(ctx, "could not retrieve schema with id ", schemaId, err.Error())
				continue
			}
			r := bytes.NewReader(msg.Value[5:])
			s3FileCreated, err := avro.DeserializeS3FileCreatedFromSchema(r, schema)
			if err != nil {
				logger.Error(ctx, "could not deserialize message ", err.Error())
				continue
			}

			// For now have the s3 downloader write to disk
			file, err := os.Create(config.OutputDirectory + s3FileCreated.Payload.Key)
			if err != nil {
				logger.Error(ctx, err)
				continue
			}
			defer file.Close()

			downloader := s3manager.NewDownloader(config.AwsSession)
			numBytes, err := downloader.Download(file,
				&s3.GetObjectInput{
					Bucket: aws.String(s3FileCreated.Payload.Bucket_name),
					Key:    aws.String(s3FileCreated.Payload.Key),
				})
			if err != nil {
				logger.Error(ctx, err)
				continue
			}

			logger.Info(ctx, "Downloaded", s3FileCreated.Payload.Key, numBytes, "bytes")
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
			operation, exists := config.OperationMap[s3FileCreated.Payload.Operation_type]
			if !exists {
				logger.Error(ctx, "invalid operation_type on file create message ")
				continue
			}
			ingestFileConfig := IngestFileConfig{
				Reader: reader,
				KafkaWriter: kafka.Writer{
					Addr:                   kafka.TCP(config.OutputBrokerAddrs...),
					Topic:                  operation.Topic,
					Logger:                 logger,
					AllowAutoTopicCreation: instrument.IsEnv("TEST"),
				},
				Tracking_id: s3FileCreated.Metadata.Tracking_id,
				Logger:      logger,
			}

			err = operation.IngestFile(ctx, ingestFileConfig)
			if err != nil {
				logger.Error(ctx, err)
				return
			}
		}
	}
}
