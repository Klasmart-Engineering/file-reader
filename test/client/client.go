package test

import (
	"context"

	"file_reader/src/log"
	"file_reader/src/protos"
	csvpb "file_reader/src/protos"
)

type Type string

var TypeName = map[Type]int32{
	"ORGANIZATION": 0,
	"SCHOOL":       1,
	"CLASS":        2,
	"USER":         3,
	"ROLE":         4,
	"PROGRAM":      5,
}

var IntToType = map[int32]Type{
	0: "ORGANIZATION",
	1: "SCHOOL",
	2: "CLASS",
	3: "USER",
	4: "ROLE",
	5: "PROGRAM",
}

type csvFileHandlers struct {
	logger *log.ZapLogger
}

type RequestBuilder struct {
	reqs []*protos.CsvFileRequest
}

func (rb RequestBuilder) getCsvFile(fileId string, filePath string, t int) *protos.CsvFileRequest {
	var typeName = protos.Type_UNKNOWN

	switch t {
	case 0:
		typeName = protos.Type_ORGANIZATION
	case 1:
		typeName = protos.Type_SCHOOL
	case 2:
		typeName = protos.Type_CLASS
	case 3:
		typeName = protos.Type_USER
	case 4:
		typeName = protos.Type_ROLE
	case 5:
		typeName = protos.Type_PROGRAM
	}
	return &protos.CsvFileRequest{
		Type:    typeName,
		Csvfile: &protos.CsvFile{FileId: fileId, Path: filePath},
	}
}
func (rb RequestBuilder) initRequests(fileIds []string, filePaths []string, typeName Type) []*protos.CsvFileRequest {

	for i := range fileIds {
		req := rb.getCsvFile(fileIds[i], filePaths[i], int(TypeName[typeName]))
		rb.reqs = append(rb.reqs, req)
	}
	return rb.reqs
}

func NewCsvFileHandlers(
	logger *log.ZapLogger,
) *csvFileHandlers {
	return &csvFileHandlers{
		logger: logger,
	}
}

func (ch *csvFileHandlers) ProcessRequests(csvClient protos.CsvFileServiceClient, req []*csvpb.CsvFileRequest) (*csvpb.CsvFileResponse, error) {
	ctx := context.Background()
	stream, err := csvClient.IngestCSV(ctx)
	if err != nil {
		ch.logger.Errorf(ctx, "Failed to get csv file: %v", err.Error())
	}

	// Iterate over the request message
	for _, v := range req {
		// Start making streaming requests by sending
		// each object inside the request message
		ch.logger.Infof(ctx, "Client streaming request: \n", v)
		stream.Send(v)
	}

	// Once the for loop finishes, the stream is closed
	// and get the response and a potential error
	res, err := stream.CloseAndRecv()
	if err != nil {
		ch.logger.Fatalf(ctx, "Error when closing the stream and receiving the response: %v", err)
	}
	return res, err
}

func (ch *csvFileHandlers) process(csvClient protos.CsvFileServiceClient, fileNames []string, filePaths []string, typeKey int32) {

	// Create a request for retrieving csv file

	typeName := IntToType[typeKey]
	req := RequestBuilder{}.initRequests(fileNames, filePaths, typeName)

	// Process request
	ch.ProcessRequests(csvClient, req)

}
