package core

import (
	"os"
	"strings"

	"github.com/KL-Engineering/file-reader/api/proto/proto_gencode/onboarding"
	"github.com/KL-Engineering/file-reader/internal/instrument"
	proto "github.com/KL-Engineering/file-reader/pkg/proto"
	protobuf "github.com/KL-Engineering/file-reader/pkg/third_party/protobuf"
)

func InitProtoOperations() Operations {
	orgTopic := instrument.MustGetEnv("ORGANIZATION_PROTO_TOPIC")
	schoolTopic := instrument.MustGetEnv("SCHOOL_PROTO_TOPIC")
	userTopic := instrument.MustGetEnv("USER_PROTO_TOPIC")
	classTopic := instrument.MustGetEnv("CLASS_PROTO_TOPIC")

	return Operations{
		OperationMap: map[string]Operation{
			"ORGANIZATION": {
				Topic:        orgTopic,
				Key:          "",
				SchemaID:     proto.SchemaRegistryClient.GetSchemaID(orgTopic),
				SerializeRow: RowToOrganizationProto,
				Headers:      OrganizationHeaders,
			},
			"SCHOOL": {
				Topic:        schoolTopic,
				Key:          "",
				SchemaID:     proto.SchemaRegistryClient.GetSchemaID(schoolTopic),
				SerializeRow: RowToSchoolProto,
				Headers:      SchoolHeaders,
			},
			"USER": {
				Topic:        userTopic,
				Key:          "",
				SchemaID:     proto.SchemaRegistryClient.GetSchemaID(userTopic),
				SerializeRow: RowToUserProto,
				Headers:      UserHeaders,
			},
			"CLASS": {
				Topic:        classTopic,
				Key:          "",
				SchemaID:     proto.SchemaRegistryClient.GetSchemaID(classTopic),
				SerializeRow: RowToClassProto,
				Headers:      ClassHeaders,
			},
		},
	}
}

func RowToOrganizationProto(row []string, tracking_id string, schemaId int, headerIndexes map[string]int) ([]byte, error) {
	md := onboarding.Metadata{
		OriginApplication: os.Getenv("METADATA_ORIGIN_APPLICATION"),
		Region:            os.Getenv("METADATA_REGION"),
		TrackingId:        tracking_id,
	}
	pl := onboarding.OrganizationPayload{
		Uuid:        row[headerIndexes[UUID]],
		Name:        row[headerIndexes[ORGANIZATION_NAME]],
		OwnerUserId: row[headerIndexes[OWNER_USER_ID]],
	}
	codec := &onboarding.Organization{Payload: &pl, Metadata: &md}
	serde := protobuf.NewProtoSerDe()
	valueBytes, err := serde.Serialize(schemaId, codec)
	if err != nil {
		return nil, err
	}
	return valueBytes, nil
}

func RowToSchoolProto(row []string, tracking_id string, schemaId int, headerIndexes map[string]int) ([]byte, error) {
	programIds := strings.Split(row[headerIndexes[PROGRAM_IDS]], ";")
	md := onboarding.Metadata{
		OriginApplication: os.Getenv("METADATA_ORIGIN_APPLICATION"),
		Region:            os.Getenv("METADATA_REGION"),
		TrackingId:        tracking_id,
	}
	pl := onboarding.SchoolPayload{
		Uuid:           &row[headerIndexes[UUID]],
		OrganizationId: row[headerIndexes[ORGANIZATION_UUID]],
		Name:           row[headerIndexes[SCHOOL_NAME]],
		ProgramIds:     programIds,
	}
	codec := &onboarding.School{Payload: &pl, Metadata: &md}
	serde := protobuf.NewProtoSerDe()
	valueBytes, err := serde.Serialize(schemaId, codec)
	if err != nil {
		return nil, err
	}
	return valueBytes, nil
}

func RowToUserProto(row []string, tracking_id string, schemaId int, headerIndexes map[string]int) ([]byte, error) {
	md := onboarding.Metadata{
		OriginApplication: os.Getenv("METADATA_ORIGIN_APPLICATION"),
		Region:            os.Getenv("METADATA_REGION"),
		TrackingId:        tracking_id,
	}
	pl := onboarding.ClassPayload{
		Uuid:             &row[headerIndexes[UUID]],
		Name:             row[headerIndexes[CLASS_NAME]],
		OrganizationUuid: row[headerIndexes[ORGANIZATION_UUID]],
	}
	codec := &onboarding.Class{Payload: &pl, Metadata: &md}
	serde := protobuf.NewProtoSerDe()
	valueBytes, err := serde.Serialize(schemaId, codec)
	if err != nil {
		return nil, err
	}
	return valueBytes, nil
}

func RowToClassProto(row []string, tracking_id string, schemaId int, headerIndexes map[string]int) ([]byte, error) {
	md := onboarding.Metadata{
		OriginApplication: os.Getenv("METADATA_ORIGIN_APPLICATION"),
		Region:            os.Getenv("METADATA_REGION"),
		TrackingId:        tracking_id,
	}
	pl := onboarding.ClassPayload{
		Uuid:             &row[headerIndexes[UUID]],
		Name:             row[headerIndexes[CLASS_NAME]],
		OrganizationUuid: row[headerIndexes[ORGANIZATION_UUID]],
	}
	codec := &onboarding.Class{Payload: &pl, Metadata: &md}
	serde := protobuf.NewProtoSerDe()
	valueBytes, err := serde.Serialize(schemaId, codec)
	if err != nil {
		return nil, err
	}
	return valueBytes, nil
}
