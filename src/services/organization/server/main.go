package main

import (
	"context"
	"file_reader/src/config"
	"file_reader/src/instrument"
	zaplogger "file_reader/src/log"
	"fmt"
	"log"
	"os"

	filepb "file_reader/src/protos/inputfile"

	src "file_reader/src"
	fileGrpc "file_reader/src/services/organization/delivery/grpc"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/riferrei/srclient"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type ingestFileServer struct {
	logger *zaplogger.ZapLogger
	cfg    *config.Config
}

func NewServer(logger *zaplogger.ZapLogger, cfg *config.Config) *ingestFileServer {
	return &ingestFileServer{
		logger: logger,
		cfg:    cfg,
	}
}
func (s *ingestFileServer) Run(ctx context.Context) error {
	s.logger.Infof(ctx, "GRPC Server is listening... at port %v\n", s.cfg.Server.Port)
	addr := instrument.GetAddressForGrpc()

	lis, grpcServer, err := instrument.GetInstrumentGrpcServer("File Processing Server", addr, s.logger)

	if err != nil {

		panic(err)
	}

	defer lis.Close()

	ingestFileService := fileGrpc.NewIngestFileService(ctx, s.logger, s.cfg)

	filepb.RegisterInputFileServiceServer(grpcServer, ingestFileService)

	s.logger.Infof(ctx, "GRPC Server is listening...", zap.String("port", s.cfg.Server.Port))

	if err = grpcServer.Serve(lis); err != nil {
		s.logger.Errorf(ctx, "Server issue.", zap.String("error", err.Error()))
	}

	return nil

}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	l, _ := zap.NewDevelopment()

	logger := zaplogger.Wrap(l)

	Logger := config.Logger{
		DisableCaller:     false,
		DisableStacktrace: false,
		Encoding:          "json",
		Level:             "info",
	}

	schemaType := os.Getenv("SCHEMA_TYPE") // AVRO or PROTO.
	if schemaType == "AVRO" {
		schemaRegistryClient := &src.SchemaRegistry{
			C: srclient.CreateSchemaRegistryClient(os.Getenv("SCHEMA_CLIENT_ENDPOINT")),
		}
		brokerAddrs := []string{os.Getenv("KAFKA_BROKER")}
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
			BrokerAddrs:     brokerAddrs,
			AwsSession:      sess,
			OutputDirectory: os.Getenv("DOWNLOAD_DIRECTORY"),
			Logger:          *log.New(os.Stdout, "kafka writer: ", 0),
		}
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokerAddrs,
			Topic:   os.Getenv("S3_FILE_CREATED_UPDATED_TOPIC"),
		})
		go src.ConsumeToIngest(ctx, r, consumerConfig, operationMap)
	}
	addr := instrument.GetAddressForGrpc()

	cfg := &config.Config{
		Server: config.Server{Port: addr, Development: true},
		Logger: Logger,
		Kafka: config.Kafka{
			Brokers: instrument.GetBrokers(),
		},
	}

	s := ingestFileServer{logger: logger, cfg: cfg}
	s.Run(ctx)
}
