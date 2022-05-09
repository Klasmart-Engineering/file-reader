package integration_test

import (
	"context"
	"file_reader/src/config"
	"file_reader/src/instrument"
	filepb "file_reader/src/protos/inputfile"
	"flag"
	"fmt"
	"net"

	"file_reader/src/log"
	fileGrpc "file_reader/src/services/organization/delivery/grpc"
	test "file_reader/test/client"

	"os"
	"testing"

	"google.golang.org/grpc/test/bufconn"

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
				InputFile: &filepb.InputFile{FileId: "file_id1", Path: "/Users/annguyen/file-reader/test/data/good/organization.csv", InputFileType: filepb.InputFileType_CSV},
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

func dialer(server *grpc.Server, service *fileGrpc.IngestFileService) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	filepb.RegisterIngestFileServiceServer(server, service)

	go func() {
		if err := server.Serve(listener); err != nil {
			fmt.Printf("Error: %v", err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}
func TestFileProcessingServer(t *testing.T) {
	flag.Set("test.timeout", "0")
	// set up env variables
	closer := envSetter(map[string]string{
		"BROKERS":          "localhost:9092",
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

	_, grpcServer, _ := instrument.GetGrpcServer(addr, logger)

	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer(grpcServer, ingestFileService)))
	if err != nil {
		fmt.Errorf(err.Error())
	}

	client := filepb.NewIngestFileServiceClient(conn)
	csvFh := test.NewInputFileHandlers(logger)

	schemaType := "PROTO"
	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			g := gomega.NewWithT(t)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// grpc call
			res, err := csvFh.ProcessRequests(ctx, client, schemaType, testCase.req)
			if testCase.name == "req ok" {

				g.Expect(res).NotTo(gomega.BeNil(), "Result should not be nil")
				g.Expect(err).To(gomega.BeNil(), "Error should be nil")
				g.Expect(res.Success).To(gomega.BeTrue())

			}

		})
	}

}
