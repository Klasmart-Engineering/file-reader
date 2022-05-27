package core

import (
	"os"

	"github.com/KL-Engineering/file-reader/api/proto/proto_gencode/onboarding"
	"github.com/KL-Engineering/file-reader/internal/instrument"
	proto "github.com/KL-Engineering/file-reader/pkg/proto"
	protobuf "github.com/KL-Engineering/file-reader/pkg/third_party/protobuf"
)

func InitProtoOperations() Operations {
	orgTopic := instrument.MustGetEnv("ORGANIZATION_PROTO_TOPIC")
	return Operations{
		OperationMap: map[string]Operation{
			"ORGANIZATION": {
				Topic:        orgTopic,
				Key:          "",
				SchemaID:     proto.SchemaRegistryClient.GetSchemaID(orgTopic),
				SerializeRow: RowToOrganizationProto,
				Headers:      OrganizationHeaders,
			},
		},
	}
}

func RowToOrganizationProto(row []string, tracking_id string, schemaId int, headerIndexes map[string]int) ([]byte, error) {
	md := onboarding.Metadata{
		OriginApplication: &onboarding.StringValue{Value: os.Getenv("METADATA_ORIGIN_APPLICATION")},
		Region:            &onboarding.StringValue{Value: os.Getenv("METADATA_REGION")},
		TrackingId:        &onboarding.StringValue{Value: tracking_id},
	}
	pl := onboarding.OrganizationPayload{
		Uuid:        &onboarding.StringValue{Value: row[headerIndexes[UUID]]},
		Name:        &onboarding.StringValue{Value: row[headerIndexes[ORGANIZATION_NAME]]},
		OwnerUserId: &onboarding.StringValue{Value: row[headerIndexes[OWNER_USER_ID]]},
	}
	codec := &onboarding.Organization{Payload: &pl, Metadata: &md}
	serde := protobuf.NewProtoSerDe()
	valueBytes, err := serde.Serialize(schemaId, codec)
	if err != nil {
		return nil, err
	}
	return valueBytes, nil
}
