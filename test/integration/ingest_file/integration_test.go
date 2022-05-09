package integration_test

import (
	"context"
	"file_reader/src/config"
	"file_reader/src/instrument"
	filepb "file_reader/src/protos/inputfile"
	"flag"
	"sync"

	"file_reader/src/log"
	fileGrpc "file_reader/src/services/organization/delivery/grpc"
	test "file_reader/test/client"
	"os"
	"testing"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"
)

var testCases = []struct {
	name        string
	req         []*filepb.InputFileRequest
	expectedRes filepb.InputFileResponse
}{
	{
		name: "req ok",
		req: []*filepb.InputFileRequest{

			&filepb.InputFileRequest{
				Type:      filepb.Type_ORGANIZATION,
				InputFile: &filepb.InputFile{FileId: "file_id1", Path: ".././test/data/good/organization.csv", InputFileType: filepb.InputFileType_CSV},
			},
		},
		expectedRes: filepb.InputFileResponse{Success: true, Errors: nil},
	},
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

func startClient(ctx context.Context, logger *log.ZapLogger, addr string, opts grpc.DialOption) filepb.InputFileServiceClient {
	con, err := grpc.Dial(addr, opts)
	if err != nil {
		logger.Fatalf(ctx, "Error connecting: %v \n", err)
	}

	defer con.Close()

	c := filepb.NewInputFileServiceClient(con)
	return c

}

func start(ctx context.Context, logger *log.ZapLogger, addr string, t *testing.T) {
	csvFh := test.NewInputFileHandlers(logger)

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
func TestFileProcessingServer(t *testing.T) {
	flag.Set("test.timeout", "0")
	// set up env variables
	closer := envSetter(map[string]string{
		"BROKERS":          "localhost:9091",
		"GRPC_SERVER":      "localhost",
		"GRPC_SERVER_PORT": "6000",
	})
	t.Cleanup(closer) // In Go 1.14+

	l, _ := zap.NewDevelopment()

	logger := log.Wrap(l)

	Logger := config.Logger{
		DisableCaller:     false,
		DisableStacktrace: false,
		Encoding:          "json",
		Level:             "info",
	}
	ctx := context.Background()
	addr := instrument.GetAddressForGrpc()

	cfg := &config.Config{
		Server: config.Server{Port: addr, Development: true},
		Logger: Logger,
		Kafka: config.Kafka{
			Brokers: instrument.GetBrokers(),
		},
	}

	ingestFileService := fileGrpc.NewIngestFileService(ctx, logger, cfg)
	lis, grpcServer, err := instrument.GetGrpcServer(addr, logger)

	if err != nil {

		t.Fatal(err)
	}

	defer lis.Close()
	t.Log(ingestFileService)
	t.Log(grpcServer)

	filepb.RegisterInputFileServiceServer(grpcServer, ingestFileService)
	// Start service client

	csvFileAddr := lis.Addr().String()
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()

	start(ctx, logger, csvFileAddr, t)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()

}
