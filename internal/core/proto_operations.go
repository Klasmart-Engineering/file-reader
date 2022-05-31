package core

import (
	"os"
	"strings"

	"github.com/KL-Engineering/file-reader/api/proto/proto_gencode/onboarding"
	"github.com/KL-Engineering/file-reader/internal/instrument"
	proto "github.com/KL-Engineering/file-reader/pkg/proto"
	protobuf "github.com/KL-Engineering/file-reader/pkg/third_party/protobuf"
)

func createRepeatedString(vals []string) []*onboarding.StringValue {
	protoStrings := []*onboarding.StringValue{}
	for _, val := range vals {
		protoStrings = append(protoStrings, &onboarding.StringValue{Value: val})
	}
	return protoStrings
}

func InitProtoOperations() Operations {
	orgTopic := instrument.MustGetEnv("ORGANIZATION_PROTO_TOPIC")
	schoolTopic := instrument.MustGetEnv("SCHOOL_PROTO_TOPIC")
	userTopic := instrument.MustGetEnv("USER_PROTO_TOPIC")
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

func RowToSchoolProto(row []string, tracking_id string, schemaId int, headerIndexes map[string]int) ([]byte, error) {
	programIds := strings.Split(row[headerIndexes[PROGRAM_IDS]], ";")
	repeatedProgramIds := createRepeatedString(programIds)
	md := onboarding.Metadata{
		OriginApplication: &onboarding.StringValue{Value: os.Getenv("METADATA_ORIGIN_APPLICATION")},
		Region:            &onboarding.StringValue{Value: os.Getenv("METADATA_REGION")},
		TrackingId:        &onboarding.StringValue{Value: tracking_id},
	}
	pl := onboarding.SchoolPayload{
		Uuid:           &onboarding.StringValue{Value: row[headerIndexes[UUID]]},
		OrganizationId: &onboarding.StringValue{Value: row[headerIndexes[ORGANIZATION_UUID]]},
		Name:           &onboarding.StringValue{Value: row[headerIndexes[SCHOOL_NAME]]},
		ProgramIds:     repeatedProgramIds,
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
		OriginApplication: &onboarding.StringValue{Value: os.Getenv("METADATA_ORIGIN_APPLICATION")},
		Region:            &onboarding.StringValue{Value: os.Getenv("METADATA_REGION")},
		TrackingId:        &onboarding.StringValue{Value: tracking_id},
	}
	pl := onboarding.UserPayload{
		Uuid:        &onboarding.StringValue{Value: row[headerIndexes[UUID]]},
		GivenName:   &onboarding.StringValue{Value: row[headerIndexes[GIVEN_NAME]]},
		FamilyName:  &onboarding.StringValue{Value: row[headerIndexes[FAMILY_NAME]]},
		Email:       &onboarding.StringValue{Value: row[headerIndexes[EMAIL]]},
		PhoneNumber: &onboarding.StringValue{Value: row[headerIndexes[PHONE_NUMBER]]},
		DateOfBirth: &onboarding.StringValue{Value: row[headerIndexes[DATE_OF_BIRTH]]},
		Gender:      &onboarding.StringValue{Value: row[headerIndexes[GENDER]]},
	}
	codec := &onboarding.User{Payload: &pl, Metadata: &md}
	serde := protobuf.NewProtoSerDe()
	valueBytes, err := serde.Serialize(schemaId, codec)
	if err != nil {
		return nil, err
	}
	return valueBytes, nil
}
