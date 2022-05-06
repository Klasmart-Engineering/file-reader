package test

import (
	"context"

	"file_reader/src/log"
	"file_reader/src/protos/csvfile"
	csvpb "file_reader/src/protos/csvfile"
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
	reqs []*csvpb.CsvFileRequest
}

func (rb RequestBuilder) getCsvFile(fileId string, filePath string, t int) *csvpb.CsvFileRequest {
	var typeName = csvpb.Type_UNKNOWN

	switch t {
	case 0:
		typeName = csvpb.Type_ORGANIZATION
	case 1:
		typeName = csvpb.Type_SCHOOL
	case 2:
		typeName = csvpb.Type_CLASS
	case 3:
		typeName = csvpb.Type_USER
	case 4:
		typeName = csvpb.Type_ROLE
	case 5:
		typeName = csvpb.Type_PROGRAM
	}
	return &csvfile.CsvFileRequest{
		Type:    typeName,
		Csvfile: &csvfile.CsvFile{FileId: fileId, Path: filePath},
	}
}
func (rb RequestBuilder) initRequests(fileIds []string, filePaths []string, typeName Type) []*csvpb.CsvFileRequest {

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

func (ch *csvFileHandlers) ProcessRequests(csvClient csvpb.CsvFileServiceClient, req []*csvpb.CsvFileRequest) (*csvpb.CsvFileResponse, error) {
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

func (ch *csvFileHandlers) process(csvClient csvfile.CsvFileServiceClient, fileNames []string, filePaths []string, typeKey int32) {

	// Create a request for retrieving csv file

	typeName := IntToType[typeKey]
	req := RequestBuilder{}.initRequests(fileNames, filePaths, typeName)

	// Process request
	ch.ProcessRequests(csvClient, req)

}
