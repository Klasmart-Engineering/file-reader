package integration_test

import (
	"bytes"
	"embed"
	"encoding/csv"
	"strings"
	"time"

	filepb "github.com/KL-Engineering/file-reader/api/proto/proto_gencode/input_file"
	"github.com/KL-Engineering/file-reader/api/proto/proto_gencode/onboarding"
	"github.com/KL-Engineering/file-reader/internal/config"
	"github.com/KL-Engineering/file-reader/internal/instrument"
	"github.com/KL-Engineering/file-reader/pkg/third_party/protobuf"
	clientPb "github.com/KL-Engineering/file-reader/test/client"
	util "github.com/KL-Engineering/file-reader/test/integration"

	"io"

	"os"
	"testing"

	"github.com/KL-Engineering/file-reader/internal/log"
	"github.com/KL-Engineering/file-reader/internal/validation"
	"github.com/KL-Engineering/file-reader/test/env"
	"go.uber.org/zap"

	"github.com/google/uuid"

	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"
	"github.com/segmentio/kafka-go"
)

//go:embed data/good
var schoolGoodDataDir embed.FS

func getSchoolCsvToProtos(filePath string) ([]*onboarding.School, error) {
	var res []*onboarding.School
	var content []byte
	content, _ = schoolGoodDataDir.ReadFile(filePath)

	// Keep store of header order
	reader := csv.NewReader(bytes.NewBuffer(content))
	headers, _ := reader.Read()
	headerIndexMap := map[string]int{}
	for i, header := range headers {
		headerIndexMap[header] = i
	}
	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				return res, nil
			}
		}
		programIds := strings.Split(row[headerIndexMap["program_ids"]], ";")

		md := onboarding.Metadata{
			OriginApplication: os.Getenv("METADATA_ORIGIN_APPLICATION"),
			Region:            os.Getenv("METADATA_REGION"),
			TrackingId:        uuid.NewString(),
		}

		pl := onboarding.SchoolPayload{
			Uuid:           &row[headerIndexMap["uuid"]],
			Name:           row[headerIndexMap["school_name"]],
			OrganizationId: row[headerIndexMap["organization_id"]],
			ProgramIds:     programIds,
		}

		res = append(res, &onboarding.School{Payload: &pl, Metadata: &md})
	}

}

func TestSchoolFileProcessingServer(t *testing.T) {
	var testCases = []struct {
		name        string
		req         []*filepb.InputFileRequest
		expectedRes filepb.InputFileResponse
	}{
		{
			name: "req ok",
			req: []*filepb.InputFileRequest{

				{
					Type:      filepb.Type_SCHOOL,
					InputFile: &filepb.InputFile{FileId: "file_id1", Path: "data/good/school.csv", InputFileType: filepb.InputFileType_CSV},
				},
			},
			expectedRes: filepb.InputFileResponse{Success: true, Errors: nil},
		},
	}

	// set up env variables
	closer := env.EnvSetter(map[string]string{
		"BROKERS":            "localhost:9092",
		"GRPC_SERVER":        "localhost",
		"GRPC_SERVER_PORT":   "6000",
		"SCHOOL_PROTO_TOPIC": uuid.NewString(),
		"SCHEMA_TYPE":        "PROTO",
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
			AllowAutoTopicCreation: instrument.IsEnv("TEST"),
		},
	}
	ctx, client := util.StartGrpc(logger, cfg, addr)

	csvFh := clientPb.NewInputFileHandlers(logger)
	schoolProtoTopic := instrument.MustGetEnv("SCHOOL_PROTO_TOPIC")

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: instrument.GetBrokers(),
		Topic:   schoolProtoTopic,
	})

	// Testing for kafka messages
	serde := protobuf.NewProtoSerDe()
	school := &onboarding.School{}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
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

				expectedValues, _ := getSchoolCsvToProtos("data/good/school.csv")
				for _, expected := range expectedValues {
					t.Log("expecting to read ", expected, " on topic ", schoolProtoTopic)
					msg, err := r.ReadMessage(ctx)
					t.Log("read message", msg, err)
					if err != nil {
						t.Logf("Error deserializing message: %v\n", err)
						break
					}

					_, err = serde.Deserialize(msg.Value, school)

					if err != nil {
						t.Logf("Error deserializing message: %v\n", err)
						break
					}

					if err == nil {
						validateTrackingId := validation.ValidateTrackingId{Uuid: school.Metadata.TrackingId}
						err := validation.UUIDValidate(validateTrackingId)
						if err != nil {
							t.Fatalf("%s", err)
						}
						g.Expect(expected.Metadata.Region).To(gomega.Equal(school.Metadata.Region))
						g.Expect(expected.Metadata.OriginApplication).To(gomega.Equal(school.Metadata.OriginApplication))
						g.Expect(util.DerefString(expected.Payload.Uuid)).To(gomega.Equal(util.DerefString(school.Payload.Uuid)))
						g.Expect(expected.Payload.Name).To(gomega.Equal(school.Payload.Name))
						g.Expect(expected.Payload.OrganizationId).To(gomega.Equal(school.Payload.OrganizationId))

					} else {
						t.Logf("Error consuming the message: %v (%v)\n", err, msg)
						break
					}
				}

			}

		})
	}
	ctx.Done()

}
