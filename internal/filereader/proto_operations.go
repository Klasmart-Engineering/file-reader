package filereader

import (
	"os"

	"github.com/KL-Engineering/file-reader/src/instrument"
	"github.com/KL-Engineering/file-reader/src/pkg/proto"
	orgPb "github.com/KL-Engineering/file-reader/src/protos/onboarding"
	"github.com/KL-Engineering/file-reader/src/third_party/protobuf"
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
			},
		},
	}
}

func RowToOrganizationProto(row []string, tracking_id string, schemaId int) ([]byte, error) {
	md := orgPb.Metadata{
		OriginApplication: &orgPb.StringValue{Value: os.Getenv("METADATA_ORIGIN_APPLICATION")},
		Region:            &orgPb.StringValue{Value: os.Getenv("METADATA_REGION")},
		TrackingId:        &orgPb.StringValue{Value: tracking_id},
	}
	pl := orgPb.OrganizationPayload{
		Uuid: &orgPb.StringValue{Value: row[0]},
		Name: &orgPb.StringValue{Value: row[1]},
	}
	codec := &orgPb.Organization{Payload: &pl, Metadata: &md}
	serde := protobuf.NewProtoSerDe()
	valueBytes, err := serde.Serialize(schemaId, codec)
	if err != nil {
		return nil, err
	}
	return valueBytes, nil
}
