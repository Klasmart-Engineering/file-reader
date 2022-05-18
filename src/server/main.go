package main

import (
	"context"
	"file_reader/src"
	"file_reader/src/config"
	"file_reader/src/instrument"
	zaplogger "file_reader/src/log"
	filepb "file_reader/src/protos/inputfile"
	"fmt"
	"os"
	"strings"

	fileGrpc "file_reader/src/services/organization/delivery/grpc"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/riferrei/srclient"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"

	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var IngestFileService *fileGrpc.IngestFileService

func grpcServerInstrument(ctx context.Context, logger *zaplogger.ZapLogger) {
	Logger := config.Logger{
		DisableCaller:     false,
		DisableStacktrace: false,
		Encoding:          "json",
		Level:             "info",
	}

	addr := instrument.GetAddressForHealthCheck()

	// grpc Server instrument
	lis, grpcServer, err := instrument.GetGrpcServer("File service health check", addr, logger)

	if err != nil {
		logger.Fatalf(ctx, "Failed to start server. Error : %v", err)
	}

	cfg := &config.Config{
		Server: config.Server{Port: addr, Development: true},
		Logger: Logger,
		Kafka: config.Kafka{
			Brokers:                instrument.GetBrokers(),
			AllowAutoTopicCreation: true,
		},
	}

	IngestFileService = fileGrpc.NewIngestFileService(ctx, logger, cfg)
	healthServer := health.NewServer()

	filepb.RegisterIngestFileServiceServer(grpcServer, IngestFileService)

	//healthService := healthcheck.NewHealthChecker()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus(filepb.IngestFileService_ServiceDesc.ServiceName, healthpb.HealthCheckResponse_SERVING)

	logger.Infof(ctx, "Server starting to listen on %s", addr)

	if err = grpcServer.Serve(lis); err != nil {
		logger.Error(ctx, "Error while starting the gRPC server on the %s listen address %v", lis, err.Error())
	}

}

func startFileCreateConsumer(ctx context.Context, logger *zaplogger.ZapLogger) {
	schemaType := os.Getenv("SCHEMA_TYPE") // AVRO or PROTO.
	if schemaType == "AVRO" {
		schemaRegistryClient := &src.SchemaRegistry{
			C:           srclient.CreateSchemaRegistryClient(os.Getenv("SCHEMA_CLIENT_ENDPOINT")),
			IdSchemaMap: make(map[int]string),
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
			fmt.Printf("Failed to initialize new aws session: %v", err)
		}
		operationMap := src.CreateOperationMap(schemaRegistryClient)
		var consumerConfig = src.ConsumeToIngestConfig{
			OutputBrokerAddrs: brokerAddrs,
			AwsSession:        sess,
			OperationMap:      operationMap,
			SchemaRegistry:    schemaRegistryClient,
			OutputDirectory:   os.Getenv("DOWNLOAD_DIRECTORY"),
			Logger:            logger,
		}
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:     brokerAddrs,
			GroupID:     os.Getenv("ORGANIZATION_GROUP_ID"),
			StartOffset: kafka.LastOffset,
			Topic:       os.Getenv("S3_FILE_CREATED_UPDATED_TOPIC"),
		})

		go src.ConsumeToIngest(ctx, r, consumerConfig)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var l *zap.Logger
	mode := instrument.MustGetEnv("MODE")
	switch mode {
	case "debug":
		l = zap.NewNop()
	default:
		l, _ = zap.NewDevelopment()
	}

	logger := zaplogger.Wrap(l)

	// Initialise Consumer for file create event
	startFileCreateConsumer(ctx, logger)

	// grpc Server instrument
	grpcServerInstrument(ctx, logger)
}
