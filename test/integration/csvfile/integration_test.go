package integration_test

import (
	"context"
	"errors"
	"file_reader/src/instrument"
	"file_reader/src/protos"
	csvpb "file_reader/src/protos"
	"sync"

	"file_reader/src/log"
	test "file_reader/test/client"
	"flag"
	"os"
	"testing"

	"google.golang.org/grpc"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"

	"go.uber.org/zap"
)

var testCases = []struct {
	name        string
	req         []*protos.CsvFileRequest
	expectedRes protos.CsvFileResponse
}{
	{
		name: "req ok",
		req: []*protos.CsvFileRequest{

			&protos.CsvFileRequest{
				Type:    protos.Type_ORGANIZATION,
				Csvfile: &protos.CsvFile{FileId: "file_id1", Path: ".././test/data/good/organization.csv"},
			},
		},
		expectedRes: protos.CsvFileResponse{Success: true, Errors: nil},
	},
}

type StreamMock struct {
	grpc.ServerStream
	ctx            context.Context
	recvToServer   chan *protos.CsvFileRequest
	sentFromClient *protos.CsvFileResponse
}

func makeStreamMock() *StreamMock {
	return &StreamMock{
		ctx:          context.Background(),
		recvToServer: make(chan *protos.CsvFileRequest, 1),
	}
}
func envSetter(envs map[string]string) (closer func()) {
	originalEnvs := map[string]string{}

	for name, value := range envs {
		if originalValue, ok := os.LookupEnv(name); ok {
			originalEnvs[name] = originalValue
		}
		_ = os.Setenv(name, value)
	}

	return func() {
		for name := range envs {
			origValue, has := originalEnvs[name]
			if has {
				_ = os.Setenv(name, origValue)
			} else {
				_ = os.Unsetenv(name)
			}
		}
	}
}

type FakeCsvFileClient struct {
	protos.UnimplementedCsvFileServiceServer
	ingestFileFn func(protos.CsvFileService_IngestCSVServer) error
}

func (fc *FakeCsvFileClient) IngestCSV(in protos.CsvFileService_IngestCSVServer) error {
	if fc.ingestFileFn != nil {
		return fc.ingestFileFn(in)
	}
	return errors.New("fakeCsvFileClient was not set up with a response - must set fc.getFileFn")
}

func startClient(ctx context.Context, logger *log.ZapLogger, addr string, opts grpc.DialOption) csvpb.CsvFileServiceClient {
	con, err := grpc.Dial(addr, opts)
	if err != nil {
		logger.Fatalf(ctx, "Error connecting: %v \n", err)
	}

	defer con.Close()

	c := csvpb.NewCsvFileServiceClient(con)
	return c

}

func start(ctx context.Context, logger *log.ZapLogger, addr string, t *testing.T) {
	csvFh := test.NewCsvFileHandlers(logger)

	opts := grpc.WithInsecure()

	c := startClient(ctx, logger, addr, opts)
	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			g := gomega.NewWithT(t)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// grpc call
			res, err := csvFh.ProcessRequests(c, testCase.req)
			if testCase.name == "req ok" {

				g.Expect(res).ToNot(gomega.BeNil(), "Result should not be nil")
				g.Expect(err).ToNot(gomega.BeNil(), "Error should be nil")
			}

		})
	}

}
func TestCsvProcessingServer(t *testing.T) {
	flag.Set("test.timeout", "0")
	// set up env variables
	closer := envSetter(map[string]string{
		"GRPC_SERVER":           "localhost",
		"GRPC_SERVER_PORT":      "6000",
		"NEW_RELIC_LICENSE_KEY": "eu01xx6802096bf8c685c39e8b3000ecFFFFNRAL",
	})
	t.Cleanup(closer) // In Go 1.14+

	l, _ := zap.NewDevelopment()

	logger := log.Wrap(l)

	addr := instrument.GetAddressForGrpc()

	fc := &FakeCsvFileClient{}
	lis, grpcServer, err := instrument.GetInstrumentGrpcServer("Testing Csv Processing Server", addr, logger)

	if err != nil {

		t.Fatal(err)
	}

	defer lis.Close()

	protos.RegisterCsvFileServiceServer(grpcServer, fc)
	// Start service client

	fakeCsvFileAddr := lis.Addr().String()
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	ctx := context.Background()
	start(ctx, logger, fakeCsvFileAddr, t)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()

}
