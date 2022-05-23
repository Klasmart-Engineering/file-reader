package grpc

import (
	"context"
	"encoding/csv"
	"file_reader/internal/filereader"
	"file_reader/src"
	"file_reader/src/config"
	"file_reader/src/instrument"
	"file_reader/src/log"
	proto "file_reader/src/pkg/proto"
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
	ctx    context.Context
	logger *log.ZapLogger
	cfg    *config.Config
	inputfile.UnimplementedIngestFileServiceServer
}

// NewIngestFileService organizationServer constructor
func NewIngestFileService(ctx context.Context, logger *log.ZapLogger, cfg *config.Config) *IngestFileService {
	return &IngestFileService{ctx: ctx, logger: logger, cfg: cfg}
}

func (c *IngestFileService) processInputFile(filePath string, fileTypeName string, schemaType string) (erroStr string) {
	// ToDo: include tracking id in proto file?
	trackingId := uuid.NewString()
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
	// Ingest file depending on schema type
	switch schemaType {
	case "AVROS":
		schemaRegistryClient := &src.SchemaRegistry{
			C: srclient.CreateSchemaRegistryClient(os.Getenv("SCHEMA_CLIENT_ENDPOINT")),
		}
		// Compose File reader for organization
		var Organization = filereader.Operation{
			Topic:        filereader.OrganizationTopicAvro,
			Key:          "",
			SchemaID:     filereader.GetOrganizationSchemaId(schemaRegistryClient),
			SerializeRow: filereader.RowToOrganizationAvro,
		}
		ingestFileConfig := filereader.IngestFileConfig{
			Reader: csv.NewReader(f),
			KafkaWriter: kafka.Writer{
				Addr:                   kafka.TCP(c.cfg.Kafka.Brokers...),
				Topic:                  filereader.OrganizationTopicAvro,
				AllowAutoTopicCreation: instrument.IsEnv("TEST"),
			},
			TrackingId: trackingId,
			Logger:     c.logger,
		}
		Organization.IngestFile(c.ctx, ingestFileConfig)
	case "PROTO":
		config := proto.Config{
			Topic:       instrument.MustGetEnv("ORGANIZATION_PROTO_TOPIC"),
			BrokerAddrs: c.cfg.Kafka.Brokers,
			Reader:      f,
			Context:     context.Background(),
			Logger:      c.logger,
		}
		errorStr := proto.OrganizationProto.IngestFilePROTO(config, fileTypeName, trackingId)
		return errorStr
	}

	return ""
}

// Ingest a new input file
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

		t := req.GetType().String()

		switch t {

		case "ORGANIZATION":

			// process organization
			if errStr := c.processInputFile(filePath, fileTypeName, "PROTO"); errStr != "" {
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

		t := req.GetType().String()

		switch t {

		case "ORGANIZATION":

			// process organization
			if errStr := c.processInputFile(filePath, fileTypeName, "AVROS"); errStr != "" {
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
}
