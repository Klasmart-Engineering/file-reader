package grpc

import (
	"context"
	"fmt"
	"io"
	"os"

	inputfile "github.com/KL-Engineering/file-reader/api/proto/proto_gencode/input_file"
	"github.com/KL-Engineering/file-reader/internal/core"

	"github.com/KL-Engineering/file-reader/internal/config"
	"github.com/KL-Engineering/file-reader/internal/instrument"
	"github.com/KL-Engineering/file-reader/internal/log"

	"github.com/google/uuid"
	"github.com/riferrei/srclient"
	"github.com/segmentio/kafka-go"
)

// ingestFileService grpc service
type IngestFileService struct {
	ctx        context.Context
	logger     *log.ZapLogger
	cfg        *config.Config
	operations core.Operations
	inputfile.UnimplementedIngestFileServiceServer
}

var operationEnumMap = map[inputfile.Type]string{
	inputfile.Type_ORGANIZATION:            "organization",
	inputfile.Type_SCHOOL:                  "school",
	inputfile.Type_CLASS:                   "class",
	inputfile.Type_USER:                    "user",
	inputfile.Type_ROLE:                    "role",
	inputfile.Type_PROGRAM:                 "program",
	inputfile.Type_ORGANIZATION_MEMBERSHIP: "organization_membership",
	inputfile.Type_CLASS_DETAILS:           "class_details",
}

// NewIngestFileService organizationServer constructor
func NewIngestFileService(ctx context.Context, logger *log.ZapLogger, cfg *config.Config) *IngestFileService {
	schemaRegistryClient := &core.SchemaRegistry{
		C:           srclient.CreateSchemaRegistryClient(os.Getenv("SCHEMA_CLIENT_ENDPOINT")),
		IdSchemaMap: make(map[int]string),
	}
	schemaType := os.Getenv("SCHEMA_TYPE") // AVRO or PROTO.
	var operations core.Operations
	switch schemaType {
	case "AVRO":
		operations = core.InitAvroOperations(schemaRegistryClient)
	case "PROTO":
		operations = core.InitProtoOperations()
	}
	return &IngestFileService{ctx: ctx, logger: logger, cfg: cfg, operations: operations}
}

func (c *IngestFileService) processInputFile(filePath string, fileTypeName string, operationType string, trackingUuid string) (erroStr string) {
	//Setup
	c.logger.Infof(c.ctx, "Processing input file in ", filePath)
	f, err := os.Open(filePath)
	if err != nil {
		c.logger.Errorf(c.ctx, "failed to open input file: ", err.Error())
		return fmt.Sprintf("%s", err)
	}
	defer f.Close()

	if err != nil {
		return fmt.Sprintf("%s", err)
	}
	// Get correct operation from operationType
	operation, exists := c.operations.GetOperation(operationType)
	if !exists {
		c.logger.Error(c.ctx, "invalid operation_type on file create message ")
	}

	fileRows := make(chan []string)
	go core.ReadRows(c.ctx, c.logger, f, "text/csv", fileRows)

	headers := <-fileRows
	headerIndexes, err := core.GetHeaderIndexes(operation.Headers, headers)

	if err != nil {
		return fmt.Sprintf("%s", err)
	}

	ingestFileConfig := core.IngestFileConfig{
		KafkaWriter: kafka.Writer{
			Addr:                   kafka.TCP(c.cfg.Kafka.Brokers...),
			Topic:                  operation.Topic,
			Logger:                 c.logger,
			AllowAutoTopicCreation: instrument.IsEnv("TEST"),
		},
		TrackingUuid: trackingUuid,
		Logger:       c.logger,
	}

	operation.IngestFile(c.ctx, fileRows, headerIndexes, ingestFileConfig)

	return ""
}

// // Ingest a new input file
func (c *IngestFileService) IngestFile(stream inputfile.IngestFileService_IngestFileServer) error {

	errors := []*inputfile.InputFileError{}

	for {
		// Start receiving stream messages from client
		req, err := stream.Recv()
		succeed := true
		if err == io.EOF {
			// Close the connection and return response to client
			if len(errors) > 0 {
				succeed = false
			}
			return stream.SendAndClose(&inputfile.InputFileResponse{Success: succeed, Errors: errors})
		}

		//Handle any possible errors when streaming requests
		if err != nil {
			c.logger.Fatalf(c.ctx, "Error when reading client request stream: %v", err)
		}

		filePath := req.InputFile.GetPath()
		fileId := req.InputFile.GetFileId()
		fileTypeName := req.InputFile.GetInputFileType().String()
		operationType := operationEnumMap[req.GetType()]
		trackingUuid := uuid.NewString()

		// process organization
		if errStr := c.processInputFile(filePath, fileTypeName, operationType, trackingUuid); errStr != "" {
			if errStr != "[]" {
				c.logger.Errorf(c.ctx, "Failed to process csv file: %s, %s", filePath, errStr)

				e := &inputfile.InputFileError{
					FileId:  fileId,
					Message: []string{"Failed to process csv file", fmt.Sprintf("Error: %s", errStr)},
				}

				// Append new error message
				errors = append(errors, e)
			}
		}
	}
}
