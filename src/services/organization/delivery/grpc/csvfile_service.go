package grpc

import (
	"context"
	"file_reader/src"
	"file_reader/src/config"
	"file_reader/src/log"
	"file_reader/src/protos/csvfile"
	"fmt"
	"io"
	"os"
)

// organizationService grpc service
type csvFileService struct {
	ctx    context.Context
	logger *log.ZapLogger
	cfg    *config.Config
	csvfile.UnimplementedCsvFileServiceServer
}

// NewOrganizationService organizationServer constructor
func NewCsvFileService(ctx context.Context, logger *log.ZapLogger, cfg *config.Config) *csvFileService {
	return &csvFileService{ctx: ctx, logger: logger, cfg: cfg}
}

func (c *csvFileService) processCsv(filePath string) error {
	//Setup
	c.logger.Infof(c.ctx, "Processing csv file in ", filePath)
	f, err := os.Open(filePath)
	if err != nil {
		c.logger.Errorf(c.ctx, "failed to open csv file: ", err.Error())
		return err
	}
	defer f.Close()

	if err != nil {
		return err
	}
	config := src.Config{
		BrokerAddrs: c.cfg.Kafka.Brokers,
		Reader:      f,
		Context:     context.Background(),
		//Logger:      ,
	}
	// Compose File reader for organization
	src.Organization.IngestFile(config)

	return nil
}

// Ingest a new csv file
func (c *csvFileService) IngestCSV(stream csvfile.CsvFileService_IngestCSVServer) error {

	errors := []*csvfile.CsvError{}
	for {
		// Start receiving stream messages from client

		req, err := stream.Recv()
		succeed := true
		if err == io.EOF {
			// Close the connection and return response to client
			if len(errors) > 0 {
				succeed = false
			}
			return stream.SendAndClose(&csvfile.CsvFileResponse{Success: succeed, Errors: errors})
		}

		//Handle any possible errors when streaming requests
		if err != nil {
			c.logger.Fatalf(c.ctx, "Error when reading client request stream: %v", err)
		}

		filePath := req.Csvfile.GetPath()
		fileId := req.Csvfile.GetFileId()
		t := req.GetType().String()

		switch t {

		case "ORGANIZATION":

			// process organization
			if err := c.processCsv(filePath); err != nil {
				c.logger.Errorf(c.ctx, "Failed to process csv file: %v, %v", filePath, err.Error())

				e := &csvfile.CsvError{
					FileId:  fileId,
					Message: []string{"Failed to process csv file", fmt.Sprint("Error: %v", err.Error())},
				}

				// Append new error message
				errors = append(errors, e)

			}

		}
	}
}
