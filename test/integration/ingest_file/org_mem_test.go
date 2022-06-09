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
	clientPb "github.com/KL-Engineering/file-reader/test/client"
	util "github.com/KL-Engineering/file-reader/test/integration"

	"github.com/KL-Engineering/file-reader/pkg/third_party/protobuf"

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
var orgMemGoodDataDir embed.FS

func getOrgMemCsvToProtos(filePath string) ([]*onboarding.OrganizationMembership, error) {
	var res []*onboarding.OrganizationMembership
	var content []byte
	content, _ = orgMemGoodDataDir.ReadFile(filePath)

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

		orgRoleUuids := strings.Split(row[headerIndexMap["organization_role_uuids"]], ";")
		md := onboarding.Metadata{
			OriginApplication: os.Getenv("METADATA_ORIGIN_APPLICATION"),
			Region:            os.Getenv("METADATA_REGION"),
			TrackingUuid:      uuid.NewString(),
		}

		pl := onboarding.OrganizationMembershipPayload{
			OrganizationUuid:      row[headerIndexMap["organization_uuid"]],
			UserUuid:              row[headerIndexMap["user_uuid"]],
			OrganizationRoleUuids: orgRoleUuids,
		}

		res = append(res, &onboarding.OrganizationMembership{Payload: &pl, Metadata: &md})
	}
}

func TestOrgMemFileProcessingServer(t *testing.T) {
	t.Skip()
	var testCases = []struct {
		name        string
		req         []*filepb.InputFileRequest
		expectedRes filepb.InputFileResponse
	}{
		{
			name: "req ok",
			req: []*filepb.InputFileRequest{

				{
					Type:      filepb.Type_ORGANIZATION_MEMBERSHIP,
					InputFile: &filepb.InputFile{FileId: "file_id1", Path: "data/good/org_mem.csv", InputFileType: filepb.InputFileType_CSV},
				},
			},
			expectedRes: filepb.InputFileResponse{Success: true, Errors: nil},
		},
	}
	// set up env variables
	closer := env.EnvSetter(map[string]string{
		"BROKERS":                             "localhost:9092",
		"GRPC_SERVER":                         "localhost",
		"GRPC_SERVER_PORT":                    "6000",
		"ORGANIZATION_MEMBERSHIP_PROTO_TOPIC": uuid.NewString(),
		"SCHEMA_TYPE":                         "PROTO",
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
	orgMemProtoTopic := instrument.MustGetEnv("ORGANIZATION_MEMBERSHIP_PROTO_TOPIC")

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: instrument.GetBrokers(),
		Topic:   orgMemProtoTopic,
	})

	// Testing for kafka messages
	serde := protobuf.NewProtoSerDe()
	orgMem := &onboarding.OrganizationMembership{}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			g := gomega.NewWithT(t)
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// grpc call
			res := csvFh.ProcessRequests(ctx, client, tc.req)
			switch tc.name {

			case "req ok":
				g.Expect(res).NotTo(gomega.BeNil(), "Result should not be nil")
				g.Expect(res.Success).To(gomega.BeTrue())

				// Testing for kafka messages
				expectedValues, _ := getOrgMemCsvToProtos("data/good/org_mem.csv")
				for _, expected := range expectedValues {
					t.Log("expecting to read ", expected, " on topic ", orgMemProtoTopic)
					msg, err := r.ReadMessage(ctx)
					t.Log("read message", msg, err)
					if err != nil {
						t.Logf("Error reading message: %v\n", err)
						t.FailNow()
					}

					_, err = serde.Deserialize(msg.Value, orgMem)

					if err != nil {
						t.Logf("Error deserializing message: %v\n", err)
						t.FailNow()
					}

					if err == nil {
						validateTrackingUuid := validation.ValidateTrackingId{Uuid: orgMem.Metadata.TrackingUuid}
						err := validation.UUIDValidate(validateTrackingUuid)
						if err != nil {
							t.Fatalf("%s", err)
						}
						g.Expect(expected.Metadata.Region).To(gomega.Equal(orgMem.Metadata.Region))
						g.Expect(expected.Metadata.OriginApplication).To(gomega.Equal(orgMem.Metadata.OriginApplication))
						g.Expect(expected.Payload.OrganizationUuid).To(gomega.Equal(orgMem.Payload.OrganizationUuid))
						g.Expect(expected.Payload.UserUuid).To(gomega.Equal(orgMem.Payload.UserUuid))
						g.Expect(expected.Payload.OrganizationRoleUuids).To(gomega.Equal(orgMem.Payload.OrganizationRoleUuids))

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
