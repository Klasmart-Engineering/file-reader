package grpc

import (
	"context"
	"file_reader/src"
	"file_reader/src/config"
	"file_reader/src/log"
	proto "file_reader/src/pkg/proto"
	"file_reader/src/protos/inputfile"
	"fmt"
	"io"
	"os"
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

func (c *IngestFileService) processInputFile(filePath string, fileTypeName string, schemaType string) error {
	//Setup
	c.logger.Infof(c.ctx, "Processing input file in ", filePath)
	f, err := os.Open(filePath)
	if err != nil {
		c.logger.Errorf(c.ctx, "failed to open input file: ", err.Error())
		return err
	}
	defer f.Close()

	if err != nil {
		return err
	}
	// Compose File reader for organization
	switch schemaType {
	case "AVROS":

		config := src.Config{
			BrokerAddrs: c.cfg.Kafka.Brokers,
			Reader:      f,
			Context:     context.Background(),
			//Logger:      ,
		}
		src.Organization.IngestFileAVROS(config, fileTypeName)
	case "PROTO":
		config := proto.Config{
			BrokerAddrs: c.cfg.Kafka.Brokers,
			Reader:      f,
			Context:     context.Background(),
			//Logger:      ,
		}
		return proto.OrganizationProto.IngestFilePROTO(config, fileTypeName)
	}

	return nil
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
			if err := c.processInputFile(filePath, fileTypeName, "PROTO"); err != nil {
				c.logger.Errorf(c.ctx, "Failed to process csv file: %v, %v", filePath, err.Error())

				e := &inputfile.InputFileError{
					FileId:  fileId,
					Message: []string{"Failed to process csv file", fmt.Sprint("Error: %v", err.Error())},
				}

				// Append new error message
				errors = append(errors, e)

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
			if err := c.processInputFile(filePath, fileTypeName, "AVROS"); err != nil {
				c.logger.Errorf(c.ctx, "Failed to process csv file: %v, %v", filePath, err.Error())

				e := &inputfile.InputFileError{
					FileId:  fileId,
					Message: []string{"Failed to process csv file", fmt.Sprint("Error: %v", err.Error())},
				}

				// Append new error message
				errors = append(errors, e)

			}

		}
	}
}
