package util

import (
	"context"
	"file_reader/internal/filereader"
	"file_reader/src"
	"file_reader/src/config"
	"file_reader/src/instrument"
	"file_reader/src/log"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	filepb "file_reader/src/protos/inputfile"
	fileGrpc "file_reader/src/services/organization/delivery/grpc"

	zaplogger "file_reader/src/log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/riferrei/srclient"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func dialer(server *grpc.Server, service *fileGrpc.IngestFileService) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	filepb.RegisterIngestFileServiceServer(server, service)

	go func() {
		if err := server.Serve(listener); err != nil {
			fmt.Printf("Error: %v", err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func StartFileCreateConsumer(ctx context.Context, logger *zaplogger.ZapLogger) {

	schemaRegistryClient := &src.SchemaRegistry{
		C:           srclient.CreateSchemaRegistryClient(os.Getenv("SCHEMA_CLIENT_ENDPOINT")),
		IdSchemaMap: make(map[int]string),
	}
	brokerAddrs := []string{"localhost:9092"}

	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "localstack",
		Config: aws.Config{
			Credentials: credentials.NewStaticCredentials(
				"test",
				"test",
				"",
			),
			Region:           aws.String("eu-west-1"),
			Endpoint:         aws.String("http://localhost:4566"),
			S3ForcePathStyle: aws.Bool(true),
		},
	})
	if err != nil {
		fmt.Printf("Failed to initialize new aws session: %v", err)
	}
	operationMap := filereader.CreateOperationMapAvro(schemaRegistryClient)
	var consumerConfig = filereader.ConsumeToIngestConfig{
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

	go filereader.ConsumeToIngest(ctx, r, consumerConfig)
}

func StartGrpc(logger *log.ZapLogger, cfg *config.Config, addr string) (context.Context, filepb.IngestFileServiceClient) {
	timeout, _ := strconv.Atoi(os.Getenv("DEFAULT_SERVER_TIMEOUT_MINS"))
	defaultTimeOut := time.Duration(timeout * int(time.Millisecond))
	ctx, _ := context.WithTimeout(context.Background(), defaultTimeOut)

	ingestFileService := fileGrpc.NewIngestFileService(ctx, logger, cfg)

	_, grpcServer, _ := instrument.GetGrpcServer("Mock service", addr, logger)

	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer(grpcServer, ingestFileService)))

	if err != nil {
		logger.Errorf(ctx, err.Error())
	}

	client := filepb.NewIngestFileServiceClient(conn)

	return ctx, client
}
