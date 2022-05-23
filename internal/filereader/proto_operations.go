package filereader

import (
	"file_reader/src/pkg/proto"
	orgPb "file_reader/src/protos/onboarding"
	"file_reader/src/third_party/protobuf"
	"os"
)

const (
	OrganizationTopicProto = "organization-proto"
)

func CreateOperationMapProto() map[string]Operation {
	// creates a map of operation_type to Operation struct
	return map[string]Operation{
		"organization": {
			Topic:        OrganizationTopicProto,
			Key:          "",
			SchemaID:     proto.SchemaRegistryClient.GetSchemaID(OrganizationTopicProto),
			SerializeRow: RowToOrganizationProto,
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
