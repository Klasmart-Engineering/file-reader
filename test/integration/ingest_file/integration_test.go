package integration_test

import (
	"bytes"
	"context"
	"embed"
	"encoding/csv"
	"file_reader/src/config"
	"file_reader/src/instrument"
	filepb "file_reader/src/protos/inputfile"

	"file_reader/src/protos/onboarding"
	orgPb "file_reader/src/protos/onboarding"
	"file_reader/src/third_party/protobuf"
	"flag"
	"fmt"
	"io"
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
	"github.com/google/uuid"
	"github.com/onsi/gomega"
	"github.com/segmentio/kafka-go"
)

//go:embed test/data/good
var testGoodDataDir embed.FS

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
func getCSVToProto(entity string) (*onboarding.Organization, error) {
	switch entity {
	case "ORGANIZATION":

		content, _ := testGoodDataDir.ReadFile("data/good/organization.csv")
		reader := csv.NewReader(bytes.NewBuffer(content))
		_, err := reader.Read() // skip first line
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
		}
		for {
			row, err := reader.Read()
			if err != nil {
				if err == io.EOF {
					fmt.Println(err)
					break
				}
			}
			md := orgPb.Metadata{
				OriginApplication: &orgPb.StringValue{Value: os.Getenv("METADATA_ORIGIN_APPLICATION")},
				Region:            &orgPb.StringValue{Value: os.Getenv("METADATA_REGION")},
				TrackingId:        &orgPb.StringValue{Value: uuid.NewString()},
			}

			pl := orgPb.OrganizationPayload{
				Uuid: &orgPb.StringValue{Value: row[0]}, //wrapperspb.String(row[0]),
				Name: &orgPb.StringValue{Value: row[1]}, //wrapperspb.String(row[1]),
			}

			return &orgPb.Organization{Payload: &pl, Metadata: &md}, nil
			fmt.Println(row)
		}

	}
	return nil
}
func TestFileProcessingServer(t *testing.T) {
	flag.Set("test.timeout", "0")
	// set up env variables
	closer := envSetter(map[string]string{
		"BROKERS":                "localhost:9092",
		"GRPC_SERVER":            "localhost",
		"GRPC_SERVER_PORT":       "6000",
		"PROTO_SCHEMA_DIRECTORY": "protos/onboarding",
	})

	defer t.Cleanup(closer) // In Go 1.14+

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
			// Testing for kafka messages
			r := kafka.NewReader(kafka.ReaderConfig{
				Brokers: instrument.GetBrokers(),
				Topic:   "organization-proto-11",
			})
			ctx := context.Background()
			serde := protobuf.NewProtoSerDe()
			org := &onboarding.Organization{}
			for i := 0; i < 5; i++ {
				msg, err := r.ReadMessage(ctx)

				if err != nil {
					fmt.Printf("Error deserializing message: %v\n", err)
				}

				_, err = serde.Deserialize(msg.Value, org)

				if err != nil {
					fmt.Printf("Error deserializing message: %v\n", err)
				}

				if err == nil {
					fmt.Printf("Received message: %v\n", org)
				} else {
					t.Logf("Error consuming the message: %v (%v)\n", err, msg)
					break
				}
			}

		})
	}

}
