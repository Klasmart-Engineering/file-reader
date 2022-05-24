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
	clientPb "file_reader/test/client"
	util "file_reader/test/integration"
	"time"

	orgPb "file_reader/src/protos/onboarding"
	"file_reader/src/third_party/protobuf"

	"io"

	"file_reader/src/log"
	"file_reader/src/pkg/validation"
	"file_reader/test/env"
	"os"
	"testing"

	"go.uber.org/zap"

	"github.com/google/uuid"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"
	"github.com/segmentio/kafka-go"
)

//go:embed data/good
var testGoodDataDir embed.FS

var testCases = []struct {
	name        string
	req         []*filepb.InputFileRequest
	expectedRes filepb.InputFileResponse
}{
	{
		name: "req ok",
		req: []*filepb.InputFileRequest{

			{
				Type:      filepb.Type_ORGANIZATION,
				InputFile: &filepb.InputFile{FileId: "file_id1", Path: "data/good/organization.csv", InputFileType: filepb.InputFileType_CSV},
			},
		},
		expectedRes: filepb.InputFileResponse{Success: true, Errors: nil},
	},
}

func getCSVToProtos(entity string, filePath string) ([]*onboarding.Organization, error) {
	var res []*onboarding.Organization
	var content []byte
	switch entity {
	case "ORGANIZATION":
		content, _ = testGoodDataDir.ReadFile(filePath)

		reader := csv.NewReader(bytes.NewBuffer(content))
		for {
			row, err := reader.Read()
			if err != nil {
				if err == io.EOF {
					return res, nil
				}
			}

			md := orgPb.Metadata{
				OriginApplication: &orgPb.StringValue{Value: os.Getenv("METADATA_ORIGIN_APPLICATION")},
				Region:            &orgPb.StringValue{Value: os.Getenv("METADATA_REGION")},
				TrackingId:        &orgPb.StringValue{Value: uuid.NewString()},
			}

			pl := orgPb.OrganizationPayload{
				Uuid: &orgPb.StringValue{Value: row[0]},
				Name: &orgPb.StringValue{Value: row[1]},
			}

			res = append(res, &orgPb.Organization{Payload: &pl, Metadata: &md})
		}

	}
	return res, nil
}

func TestFileProcessingServer(t *testing.T) {
	// set up env variables
	closer := env.EnvSetter(map[string]string{
		"BROKERS":                  "localhost:9092",
		"GRPC_SERVER":              "localhost",
		"GRPC_SERVER_PORT":         "6000",
		"ORGANIZATION_PROTO_TOPIC": uuid.NewString(),
		"SCHEMA_TYPE":              "PROTO",
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
	addr := instrument.GetAddressForGrpc()

	cfg := &config.Config{
		Server: config.Server{Port: addr, Development: true},
		Logger: Logger,
		Kafka: config.Kafka{
			Brokers:                instrument.GetBrokers(),
			DialTimeout:            int(3 * time.Minute),
			MaxAttempts:            3,
			AllowAutoTopicCreation: true,
		},
	}
	ctx, client := util.StartGrpc(logger, cfg, addr)

	csvFh := clientPb.NewInputFileHandlers(logger)
	orgProtoTopic := instrument.MustGetEnv("ORGANIZATION_PROTO_TOPIC")
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: instrument.GetBrokers(),
		Topic:   orgProtoTopic,
	})

	// Testing for kafka messages
	serde := protobuf.NewProtoSerDe()
	org := &onboarding.Organization{}

	for _, tc := range testCases {
		testCase := tc

		t.Run(tc.name, func(t *testing.T) {
			g := gomega.NewWithT(t)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// grpc call
			res := csvFh.ProcessRequests(ctx, client, testCase.req)
			switch testCase.name {

			case "req ok":
				g.Expect(res).NotTo(gomega.BeNil(), "Result should not be nil")
				g.Expect(res.Success).To(gomega.BeTrue())

				// Testing for kafka messages

				ctx := context.Background()

				expectedValues, _ := getCSVToProtos("ORGANIZATION", "data/good/organization.csv")
				for _, expected := range expectedValues {
					t.Log("expecting to read ", expected, " on topic ", orgProtoTopic)
					msg, err := r.ReadMessage(ctx)
					t.Log("read message", msg, err)
					if err != nil {
						t.Logf("Error deserializing message: %v\n", err)
						break
					}

					_, err = serde.Deserialize(msg.Value, org)

					if err != nil {
						t.Logf("Error deserializing message: %v\n", err)
						break
					}

					if err == nil {
						validateTrackingId := validation.ValidateTrackingId{Uuid: org.Metadata.TrackingId.Value}
						err := validation.UUIDValidate(validateTrackingId)
						if err != nil {
							t.Fatalf("%s", err)
						}
						g.Expect(expected.Metadata.Region.Value).To(gomega.Equal(org.Metadata.Region.Value))
						g.Expect(expected.Metadata.OriginApplication.Value).To(gomega.Equal(org.Metadata.OriginApplication.Value))
						g.Expect(expected.Payload.Uuid.Value).To(gomega.Equal(org.Payload.Uuid.Value))
						g.Expect(expected.Payload.Name.Value).To(gomega.Equal(org.Payload.Name.Value))

					} else {
						t.Logf("Error consuming the message: %v (%v)\n", err, msg)
						break
					}
				}

			}

		})
	}

}
