package filereader

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/csv"
	"io/ioutil"
	"log"
	"os"
	"strings"

	avro "github.com/KL-Engineering/file-reader/avro_gencode"
	"github.com/KL-Engineering/file-reader/src"
	"github.com/KL-Engineering/file-reader/src/instrument"
	zaplogger "github.com/KL-Engineering/file-reader/src/log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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
	SchemaRegistry    *src.SchemaRegistry
	Operations        Operations
	Logger            *zaplogger.ZapLogger
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

			// Create and open file on /tmp/
			f, err := ioutil.TempFile("", "file-reader-"+s3FileCreated.Payload.Key)
			if err != nil {
				log.Fatal("Failed to make tmp file", err)
			}
			defer os.Remove(f.Name())

			// Download from S3 to file
			downloader := s3manager.NewDownloader(config.AwsSession)
			numBytes, err := downloader.Download(f,
				&s3.GetObjectInput{
					Bucket: aws.String(s3FileCreated.Payload.Bucket_name),
					Key:    aws.String(s3FileCreated.Payload.Key),
				})
			if err != nil {
				logger.Error(ctx, err)
				continue
			}
			logger.Infof(ctx, "Downloaded %s %d bytes", s3FileCreated.Payload.Key, numBytes)

			// Close and reopen the same file for ingest (until thought of alternative)
			f.Close()
			f, _ = os.Open(f.Name())
			defer f.Close()
			// Compose the ingestFile() with different reader depending on file type
			var reader Reader
			switch s3FileCreated.Payload.Content_type {
			default:
				reader = csv.NewReader(f)
			}

			// Map to operation based on operation type
			operation, exists := config.Operations.GetOperation(s3FileCreated.Payload.Operation_type)
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
				TrackingId: s3FileCreated.Metadata.Tracking_id,
				Logger:     logger,
			}

			operation.IngestFile(ctx, ingestFileConfig)
			f.Close()
		}
	}
}

func StartFileCreateConsumer(ctx context.Context, logger *zaplogger.ZapLogger) {
	schemaRegistryClient := &src.SchemaRegistry{
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
		Profile: os.Getenv("AWS_PROFILE"),
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
	}
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokerAddrs,
		GroupID:     os.Getenv("S3_FILE_CREATED_UPDATED_GROUP_ID"),
		StartOffset: kafka.LastOffset,
		Topic:       os.Getenv("S3_FILE_CREATED_UPDATED_TOPIC"),
	})

	go ConsumeToIngest(ctx, r, consumerConfig)
}
