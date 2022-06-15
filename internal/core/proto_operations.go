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
	orgMemTopic := instrument.MustGetEnv("ORGANIZATION_MEMBERSHIP_PROTO_TOPIC")
	classDetailsTopic := instrument.MustGetEnv("CLASS_DETAILS_PROTO_TOPIC")
	schoolMemTopic := instrument.MustGetEnv("SCHOOL_MEMBERSHIP_PROTO_TOPIC")

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

			"ORGANIZATION_MEMBERSHIP": {
				Topic:        orgMemTopic,
				Key:          "",
				SchemaID:     proto.SchemaRegistryClient.GetSchemaID(orgMemTopic),
				SerializeRow: RowToOrgMemProto,
				Headers:      OrgMemHeaders,
			},
			"CLASS_DETAILS": {
				Topic:        classDetailsTopic,
				Key:          "",
				SchemaID:     proto.SchemaRegistryClient.GetSchemaID(classDetailsTopic),
				SerializeRow: RowToClassDetailsProto,
				Headers:      ClassDetailsHeaders,
			},
			"SCHOOL_MEMBERSHIP": {
				Topic:        schoolMemTopic,
				Key:          "",
				SchemaID:     proto.SchemaRegistryClient.GetSchemaID(schoolMemTopic),
				SerializeRow: RowToSchoolMemProto,
				Headers:      SchoolMemHeaders,
			},
		},
	}
}

func RowToOrganizationProto(row []string, tracking_uuid string, schemaId int, headerIndexes map[string]int) ([]byte, error) {
	md := onboarding.Metadata{
		OriginApplication: os.Getenv("METADATA_ORIGIN_APPLICATION"),
		Region:            os.Getenv("METADATA_REGION"),
		TrackingUuid:      tracking_uuid,
	}
	pl := onboarding.OrganizationPayload{
		Uuid:          row[headerIndexes[UUID]],
		Name:          row[headerIndexes[NAME]],
		OwnerUserUuid: row[headerIndexes[OWNER_USER_UUID]],
	}
	codec := &onboarding.Organization{Payload: &pl, Metadata: &md}
	serde := protobuf.NewProtoSerDe()
	valueBytes, err := serde.Serialize(schemaId, codec)
	if err != nil {
		return nil, err
	}
	return valueBytes, nil
}

func RowToOrgMemProto(row []string, tracking_uuid string, schemaId int, headerIndexes map[string]int) ([]byte, error) {
	orgRoleUuids := strings.Split(row[headerIndexes[ORGANIZATION_ROLE_UUIDS]], ";")
	md := onboarding.Metadata{
		OriginApplication: os.Getenv("METADATA_ORIGIN_APPLICATION"),
		Region:            os.Getenv("METADATA_REGION"),
		TrackingUuid:      tracking_uuid,
	}
	pl := onboarding.OrganizationMembershipPayload{
		OrganizationUuid:      row[headerIndexes[ORGANIZATION_UUID]],
		UserUuid:              row[headerIndexes[USER_UUID]],
		OrganizationRoleUuids: orgRoleUuids,
	}
	codec := &onboarding.OrganizationMembership{Payload: &pl, Metadata: &md}
	serde := protobuf.NewProtoSerDe()
	valueBytes, err := serde.Serialize(schemaId, codec)
	if err != nil {
		return nil, err
	}
	return valueBytes, nil
}

func RowToSchoolProto(row []string, tracking_uuid string, schemaId int, headerIndexes map[string]int) ([]byte, error) {
	programUuids := strings.Split(row[headerIndexes[PROGRAM_UUIDS]], ";")
	md := onboarding.Metadata{
		OriginApplication: os.Getenv("METADATA_ORIGIN_APPLICATION"),
		Region:            os.Getenv("METADATA_REGION"),
		TrackingUuid:      tracking_uuid,
	}
	pl := onboarding.SchoolPayload{
		Uuid:             &row[headerIndexes[UUID]],
		OrganizationUuid: row[headerIndexes[ORGANIZATION_UUID]],
		Name:             row[headerIndexes[NAME]],
		ProgramUuids:     programUuids,
	}
	codec := &onboarding.School{Payload: &pl, Metadata: &md}
	serde := protobuf.NewProtoSerDe()
	valueBytes, err := serde.Serialize(schemaId, codec)
	if err != nil {
		return nil, err
	}
	return valueBytes, nil
}

func RowToSchoolMemProto(row []string, tracking_uuid string, schemaId int, headerIndexes map[string]int) ([]byte, error) {
	md := onboarding.Metadata{
		OriginApplication: os.Getenv("METADATA_ORIGIN_APPLICATION"),
		Region:            os.Getenv("METADATA_REGION"),
		TrackingUuid:      tracking_uuid,
	}
	pl := onboarding.SchoolMembershipPayload{
		SchoolUuid: row[headerIndexes[SCHOOL_UUID]],
		UserUuid:   row[headerIndexes[USER_UUID]],
	}
	codec := &onboarding.SchoolMembership{Payload: &pl, Metadata: &md}
	serde := protobuf.NewProtoSerDe()
	valueBytes, err := serde.Serialize(schemaId, codec)
	if err != nil {
		return nil, err
	}
	return valueBytes, nil
}

func RowToUserProto(row []string, tracking_uuid string, schemaId int, headerIndexes map[string]int) ([]byte, error) {
	md := onboarding.Metadata{
		OriginApplication: os.Getenv("METADATA_ORIGIN_APPLICATION"),
		Region:            os.Getenv("METADATA_REGION"),
		TrackingUuid:      tracking_uuid,
	}
	pl := onboarding.UserPayload{
		Uuid:        row[headerIndexes[UUID]],
		GivenName:   row[headerIndexes[GIVEN_NAME]],
		FamilyName:  row[headerIndexes[FAMILY_NAME]],
		Email:       &row[headerIndexes[EMAIL]],
		PhoneNumber: &row[headerIndexes[PHONE_NUMBER]],
		DateOfBirth: &row[headerIndexes[DATE_OF_BIRTH]],
		Gender:      row[headerIndexes[GENDER]],
	}
	codec := &onboarding.User{Payload: &pl, Metadata: &md}
	serde := protobuf.NewProtoSerDe()
	valueBytes, err := serde.Serialize(schemaId, codec)
	if err != nil {
		return nil, err
	}
	return valueBytes, nil
}

func RowToClassProto(row []string, tracking_uuid string, schemaId int, headerIndexes map[string]int) ([]byte, error) {
	md := onboarding.Metadata{
		OriginApplication: os.Getenv("METADATA_ORIGIN_APPLICATION"),
		Region:            os.Getenv("METADATA_REGION"),
		TrackingUuid:      tracking_uuid,
	}
	pl := onboarding.ClassPayload{
		Uuid:             &row[headerIndexes[UUID]],
		Name:             row[headerIndexes[NAME]],
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

func RowToClassDetailsProto(row []string, tracking_uuid string, schemaId int, headerIndexes map[string]int) ([]byte, error) {
	md := onboarding.Metadata{
		OriginApplication: os.Getenv("METADATA_ORIGIN_APPLICATION"),
		Region:            os.Getenv("METADATA_REGION"),
		TrackingUuid:      tracking_uuid,
	}
	pl := onboarding.ClassDetailsPayload{
		ClassUuid:        row[headerIndexes[CLASS_UUID]],
		SchoolUuid:       &row[headerIndexes[SCHOOL_UUID]],
		ProgramUuids:     strings.Split(row[headerIndexes[PROGRAM_UUIDS]], ";"),
		SubjectUuids:     strings.Split(row[headerIndexes[SUBJECT_UUIDS]], ";"),
		GradeUuids:       strings.Split(row[headerIndexes[GRADE_UUIDS]], ";"),
		AgeRangeUuids:    strings.Split(row[headerIndexes[AGE_RANGE_UUIDS]], ";"),
		AcademicTermUuid: &row[headerIndexes[ACADEMIC_TERM_UUID]],
	}
	codec := &onboarding.ClassDetails{Payload: &pl, Metadata: &md}
	serde := protobuf.NewProtoSerDe()
	valueBytes, err := serde.Serialize(schemaId, codec)
	if err != nil {
		return nil, err
	}
	return valueBytes, nil
}
