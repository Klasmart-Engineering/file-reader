package test

import (
	"context"
	"time"

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

type organizationHandlers struct {
	log            zap.Logger
	organizationUC organization.UseCase
	validate       *validator.Validate
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
	return &o
}

func (ch *csvFileHandlers) ProcessRequests(csvClient protos.CsvFileServiceClient, req []*csvpb.CsvFileRequest) (*csvpb.CsvFileResponse, error) {
	ctx := context.Background()
	stream, err := csvClient.IngestCSV(ctx)
	if err != nil {
		ch.logger.Errorf(ctx, "Failed to get csv file: %v", err.Error())
	}
	fmt.Printf("Schema '%d' retrieved successfully!\n", schema.ID())

	return func(f *os.File, ctx context.Context) {
		csvReader := csv.NewReader(f)
		for {
			row, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				oh.log.Sugar().Fatalf(err.Error())
			}

			// Map row to bytes using schema
			op := rowToOpConverter(row)
			if err := oh.organizationUC.PublishCreate(ctx, op, schema); err != nil {
				oh.log.Sugar().Errorf("organizationUC.PublishCreate: %v", err)
			}
		}
	}
}

func CreateOrganization(ctx context.Context, schemaText string, f *os.File, log zap.Logger) {

		ch.logger.Infof(ctx, "Client streaming request: %v\n", v)
		stream.Send(v)
		time.Sleep(500 * time.Millisecond)
	}

	// Once the all requests are received, the stream is closed
	// get the response and potential errors
	res, err := stream.CloseAndRecv()
	if err != nil {
		ch.logger.Errorf(ctx, "Error when closing the stream and receiving the response: %v\n", err)
	}
	orgProducer := kafka.NewOrganizationProducer(log, cfg)
	// Initialize producer writer
	orgProducer.Run()
	organizationUC := usecase.NewOrganizationUC(log, orgProducer)
	handler := NewOrganizationHandlers(log, organizationUC, validate)

	if len(res.Errors) > 0 {
		ch.logger.Errorf(ctx, "Csv processing error: %v\n", res.Errors)
		return res, err
	}

	return res, nil
}
func (ch *csvFileHandlers) process(csvClient protos.CsvFileServiceClient, fileNames []string, filePaths []string, typeKey int32) {

	// Create a request for retrieving csv file

	typeName := IntToType[typeKey]
	req := RequestBuilder{}.initRequests(fileNames, filePaths, typeName)

	// Process request
	ch.ProcessRequests(csvClient, req)

}
