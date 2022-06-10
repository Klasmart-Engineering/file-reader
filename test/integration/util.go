package util

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/KL-Engineering/file-reader/internal/config"
	"github.com/KL-Engineering/file-reader/internal/core"
	"github.com/KL-Engineering/file-reader/internal/instrument"
	zapLogger "github.com/KL-Engineering/file-reader/internal/log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/riferrei/srclient"
	"github.com/segmentio/kafka-go"

	avro "github.com/KL-Engineering/file-reader/api/avro/avro_gencode"
	filepb "github.com/KL-Engineering/file-reader/api/proto/proto_gencode/input_file"
	fileGrpc "github.com/KL-Engineering/file-reader/internal/services/delivery/grpc"

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

func StartGrpc(logger *zapLogger.ZapLogger, cfg *config.Config, addr string) (context.Context, filepb.IngestFileServiceClient, net.Listener) {

	timeout, _ := strconv.Atoi(os.Getenv("DEFAULT_SERVER_TIMEOUT_MS"))
	defaultTimeOut := time.Duration(timeout * int(time.Millisecond))
	ctx, _ := context.WithTimeout(context.Background(), defaultTimeOut)

	ingestFileService := fileGrpc.NewIngestFileService(ctx, logger, cfg)

	ln, grpcServer, _ := instrument.GetGrpcServer("Mock service", addr, logger)
	conn, err := grpc.DialContext(ctx, "", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(dialer(grpcServer, ingestFileService)))

	if err != nil {
		logger.Errorf(ctx, err.Error())
	}

	client := filepb.NewIngestFileServiceClient(conn)

	return ctx, client, ln
}

func UploadFileToS3(bucket string, s3key string, awsRegion string, file io.Reader) error {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Credentials: credentials.NewStaticCredentials(
				"test",
				"test",
				"",
			),
			Region:           aws.String("awsRegion"),
			Endpoint:         aws.String("http://localhost:4566"),
			S3ForcePathStyle: aws.Bool(true),
		},
	})
	if err != nil {
		return err
	}
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(s3key),
		Body:   file,
	})
	if err != nil {
		return err
	}
	return nil
}

func ProduceFileCreateMessage(ctx context.Context, s3FileCreationTopic string, brokerAddrs []string, s3FileCreated avro.S3FileCreatedUpdated) error {
	schemaRegistryClient := &core.SchemaRegistry{
		C: srclient.CreateSchemaRegistryClient("http://localhost:8081"),
	}
	schemaBody := avro.S3FileCreatedUpdated.Schema(avro.NewS3FileCreatedUpdated())
	s3FileCreationSchemaId := schemaRegistryClient.GetSchemaId(schemaBody, s3FileCreationTopic)
	schemaIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(schemaIDBytes, uint32(s3FileCreationSchemaId))

	var buf bytes.Buffer
	s3FileCreated.Serialize(&buf)
	valueBytes := buf.Bytes()

	// Combine row bytes with schema id to make a record
	var recordValue []byte
	recordValue = append(recordValue, byte(0))
	recordValue = append(recordValue, schemaIDBytes...)
	recordValue = append(recordValue, valueBytes...)

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
		return err
	}
	return nil
}

func DerefString(s *string) string {
	if s != nil {
		return *s
	}

	return ""
}

func DerefAvroNullString(s *avro.UnionNullString) string {
	if s != nil {
		return s.String
	}

	return ""
}
