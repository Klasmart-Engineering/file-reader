package grpc

import (
	"context"
	"encoding/csv"
	"file_reader/internal/filereader"
	"file_reader/src"
	"file_reader/src/config"
	"file_reader/src/instrument"
	"file_reader/src/log"
	"file_reader/src/protos/inputfile"
	"fmt"
	"io"
	"os"

	"github.com/google/uuid"
	"github.com/riferrei/srclient"
	"github.com/segmentio/kafka-go"
)

// ingestFileService grpc service
type IngestFileService struct {
	ctx        context.Context
	logger     *log.ZapLogger
	cfg        *config.Config
	operations filereader.Operations
	inputfile.UnimplementedIngestFileServiceServer
}

var operationEnumMap = map[inputfile.Type]string{
	inputfile.Type_ORGANIZATION: "organization",
	inputfile.Type_SCHOOL:       "school",
	inputfile.Type_CLASS:        "class",
	inputfile.Type_USER:         "user",
	inputfile.Type_ROLE:         "role",
	inputfile.Type_PROGRAM:      "program",
}

// NewIngestFileService organizationServer constructor
func NewIngestFileService(ctx context.Context, logger *log.ZapLogger, cfg *config.Config) *IngestFileService {
	schemaRegistryClient := &src.SchemaRegistry{
		C:           srclient.CreateSchemaRegistryClient(os.Getenv("SCHEMA_CLIENT_ENDPOINT")),
		IdSchemaMap: make(map[int]string),
	}
	schemaType := os.Getenv("SCHEMA_TYPE") // AVRO or PROTO.
	var operations filereader.Operations
	switch schemaType {
	case "AVRO":
		operations = filereader.InitAvroOperations(schemaRegistryClient)
	case "PROTO":
		operations = filereader.InitProtoOperations()
	}
	return &IngestFileService{ctx: ctx, logger: logger, cfg: cfg, operations: operations}
}

func (c *IngestFileService) processInputFile(filePath string, fileTypeName string, schemaType string, operationType string, trackingId string) (erroStr string) {
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

	ingestFileConfig := filereader.IngestFileConfig{
		Reader: csv.NewReader(f),
		KafkaWriter: kafka.Writer{
			Addr:                   kafka.TCP(c.cfg.Kafka.Brokers...),
			Topic:                  operation.Topic,
			Logger:                 c.logger,
			AllowAutoTopicCreation: instrument.IsEnv("TEST"),
		},
		TrackingId: trackingId,
		Logger:     c.logger,
	}

	operation.IngestFile(c.ctx, ingestFileConfig)

	return ""
}

// // Ingest a new input file
func (c *IngestFileService) IngestFilePROTO(stream inputfile.IngestFileService_IngestFilePROTOServer) error {

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
		trackingId := uuid.NewString()

		t := req.GetType().String()

		switch t {

		case "ORGANIZATION":

			// process organization
			if errStr := c.processInputFile(filePath, fileTypeName, "PROTO", operationType, trackingId); errStr != "" {
				if errStr != "[]" {
					c.logger.Errorf(c.ctx, "Failed to process csv file: %s, %s", filePath, errStr)

					e := &inputfile.InputFileError{
						FileId:  fileId,
						Message: []string{"Failed to process csv file", fmt.Sprint("Error: %s", errStr)},
					}

					// Append new error message
					errors = append(errors, e)
				}

			}

		}
	}
}

// Ingest a new input file
func (c *IngestFileService) IngestFileAVROS(stream inputfile.IngestFileService_IngestFileAVROSServer) error {

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
		trackingId := uuid.NewString()

		if errStr := c.processInputFile(filePath, fileTypeName, "AVROS", operationType, trackingId); errStr != "" {
			c.logger.Errorf(c.ctx, "Failed to process input file: %s, %s", filePath, errStr)

			e := &inputfile.InputFileError{
				FileId:  fileId,
				Message: []string{"Failed to process input file", fmt.Sprint("Error: %s", errStr)},
			}

			// Append new error message
			errors = append(errors, e)

		}

	}

}
