package core

import (
	"context"
	"os"
	"strings"

	"github.com/KL-Engineering/file-reader/internal/instrument"
	zaplogger "github.com/KL-Engineering/file-reader/internal/log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/riferrei/srclient"
	"github.com/segmentio/kafka-go"
)

type Operations struct {
	OperationMap map[string]Operation
}

func (ops Operations) GetOperation(opKey string) (Operation, bool) {
	op, exists := ops.OperationMap[strings.ToUpper(opKey)]
	return op, exists
}

type ConsumeToIngestConfig struct {
	OutputBrokerAddrs []string
	AwsSession        *session.Session
	SchemaRegistry    *SchemaRegistry
	Operations        Operations
	Logger            *zaplogger.ZapLogger
	BatchSize         int
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
			logger.Debug(ctx, " received message: ", string(msg.Value))

			// Deserialize file create message
			s3FileCreated, err := deserializeS3Event(config.SchemaRegistry, msg.Value)
			if err != nil {
				logger.Error(ctx, "could not deserialize message ", err.Error())
				continue
			}

			// Download S3 file
			fileRows := make(chan []string)
			go StreamFile(ctx, config.Logger, config.AwsSession, s3FileCreated, fileRows, config.BatchSize)

			// Map to operation based on operation type
			operation, exists := config.Operations.GetOperation(s3FileCreated.Payload.Operation_type)

			if !exists {
				logger.Error(ctx, "invalid operation_type on file create message ")
				continue
			}
			// Parse file headers
			headers := <-fileRows
			headerIndexes, err := GetHeaderIndexes(operation.Headers, headers)
			if err != nil {
				logger.Error(ctx, "Error parsing file headers. ", err.Error())
				continue
			}

			ingestFileConfig := IngestFileConfig{
				KafkaWriter: kafka.Writer{
					Addr:                   kafka.TCP(config.OutputBrokerAddrs...),
					Topic:                  operation.Topic,
					Logger:                 logger,
					AllowAutoTopicCreation: instrument.IsEnv("TEST"),
				},
				TrackingUuid: s3FileCreated.Metadata.Tracking_uuid,
				Logger:       logger,
			}

			operation.IngestFile(ctx, fileRows, headerIndexes, ingestFileConfig)
		}
	}
}

func StartFileCreateConsumer(ctx context.Context, logger *zaplogger.ZapLogger) {
	schemaRegistryClient := &SchemaRegistry{
		C:           srclient.CreateSchemaRegistryClient(os.Getenv("SCHEMA_CLIENT_ENDPOINT")),
		IdSchemaMap: make(map[int]string),
	}

	schemaType := os.Getenv("SCHEMA_TYPE") // AVRO or PROTO.
	var operations Operations
	switch schemaType {
	case "AVRO":
		operations = InitAvroOperations(schemaRegistryClient)
	case "PROTO":
		operations = InitProtoOperations()
	}

	brokerAddrs := strings.Split(os.Getenv("BROKERS"), ",")
	sess, err := session.NewSessionWithOptions(session.Options{
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
		logger.Infof(ctx, "Failed to initialize new aws session: %s", err)
	}

	var consumerConfig = ConsumeToIngestConfig{
		OutputBrokerAddrs: brokerAddrs,
		AwsSession:        sess,
		Operations:        operations,
		SchemaRegistry:    schemaRegistryClient,
		Logger:            logger,
		BatchSize:         100000,
	}
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokerAddrs,
		GroupID:     os.Getenv("S3_FILE_CREATED_UPDATED_GROUP_ID"),
		StartOffset: kafka.FirstOffset,
		Topic:       os.Getenv("S3_FILE_CREATED_UPDATED_TOPIC"),
	})

	go ConsumeToIngest(ctx, r, consumerConfig)
}
