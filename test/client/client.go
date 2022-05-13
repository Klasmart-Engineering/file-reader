package client

import (
	"context"
	"file_reader/src/log"
	filepb "file_reader/src/protos/inputfile"
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

type InputFileType string

var InputFileTypeName = map[InputFileType]int32{
	"CSV": 0,
}

var IntToFileType = map[int32]InputFileType{
	0: "CSV",
}

type inputFileHandlers struct {
	logger *log.ZapLogger
}

type RequestBuilder struct {
	reqs []*filepb.InputFileRequest
}

func (rb RequestBuilder) getInputFile(fileId string, filePath string, entity int, fileType int) *filepb.InputFileRequest {
	var typeName = filepb.Type_UNKNOWN
	var fileTypeName = filepb.InputFileType_CSV

	switch entity {
	case 0:
		typeName = filepb.Type_ORGANIZATION
	case 1:
		typeName = filepb.Type_SCHOOL
	case 2:
		typeName = filepb.Type_CLASS
	case 3:
		typeName = filepb.Type_USER
	case 4:
		typeName = filepb.Type_ROLE
	case 5:
		typeName = filepb.Type_PROGRAM
	}

	switch fileType {
	case 0:
		fileTypeName = filepb.InputFileType_CSV
	}

	return &filepb.InputFileRequest{
		Type:      typeName,
		InputFile: &filepb.InputFile{FileId: fileId, InputFileType: fileTypeName, Path: filePath},
	}
}
func (rb RequestBuilder) initRequests(fileIds []string, filePaths []string, entityTypeName Type, fileTypeName InputFileType) []*filepb.InputFileRequest {

	for i := range fileIds {
		req := rb.getInputFile(fileIds[i], filePaths[i], int(TypeName[entityTypeName]), int(InputFileTypeName[fileTypeName]))
		rb.reqs = append(rb.reqs, req)
	}
	return rb.reqs
}

func NewInputFileHandlers(
	logger *log.ZapLogger,
) *inputFileHandlers {
	return &inputFileHandlers{
		logger: logger,
	}
}

func (h *inputFileHandlers) ProcessRequests(ctx context.Context, fileClient filepb.IngestFileServiceClient, schemaType string, req []*filepb.InputFileRequest) (*filepb.InputFileResponse, error) {
	var stream filepb.IngestFileService_IngestFilePROTOClient
	var err error

	switch schemaType {
	case "AVROS":
		stream, err = fileClient.IngestFileAVROS(ctx)
	case "PROTO":
		stream, err = fileClient.IngestFilePROTO(ctx)
	}
	if err != nil {
		h.logger.Errorf(ctx, false, "Error on IngestFile rpc call: %v", err.Error())
	}

	// Iterate over the request message
	for _, v := range req {
		// Start making streaming requests by sending
		// each object inside the request message
		h.logger.Infof(ctx, false, "Client streaming request: \n", v)
		stream.Send(v)
	}

	// Once the for loop finishes, the stream is closed
	// and get the response and a potential error
	res, err := stream.CloseAndRecv()
	if err != nil {
		h.logger.Fatalf(ctx, false, "Error when closing the stream and receiving the response: %v", err)
	}
	return res, err
}

func (h *inputFileHandlers) process(ctx context.Context, fileClient filepb.IngestFileServiceClient, schemaType string, fileNames []string, filePaths []string, entityTypeKey int32, inputFileTypeKey int32) (*filepb.InputFileResponse, error) {

	// Create a request for retrieving csv file

	entityTypeName := IntToType[entityTypeKey]
	inputFileTypeName := IntToFileType[inputFileTypeKey]

	req := RequestBuilder{}.initRequests(fileNames, filePaths, entityTypeName, inputFileTypeName)
	// Process request
	return h.ProcessRequests(ctx, fileClient, schemaType, req)

}
