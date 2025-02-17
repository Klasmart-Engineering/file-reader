package core

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"strings"

	avrogen "github.com/KL-Engineering/file-reader/api/avro/avro_gencode"
	"github.com/KL-Engineering/file-reader/internal/instrument"
)

type avroCodec interface {
	// Represents the structs in the avro-gencode which have Serialize() functions
	Serialize(io.Writer) error
}

func makeAvroOptionalString(value string) *avrogen.UnionNullString {
	if value != "" {
		return &avrogen.UnionNullString{
			String:    value,
			UnionType: avrogen.UnionNullStringTypeEnumString,
		}
	}
	return nil
}

func makeAvroOptionalArrayString(value string) *avrogen.UnionNullArrayString {
	if value != "" {
		return &avrogen.UnionNullArrayString{
			ArrayString: strings.Split(value, ";"),
			UnionType:   avrogen.UnionNullArrayStringTypeEnumArrayString,
		}
	}
	return nil
}

func serializeAvroRecord(codec avroCodec, schemaId int) []byte {
	// Get bytes for the schemaId
	schemaIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(schemaIDBytes, uint32(schemaId))

	// Get bytes for the row
	var buf bytes.Buffer
	codec.Serialize(&buf)
	valueBytes := buf.Bytes()

	//Combine row bytes with schema id to make a record
	var recordValue []byte
	recordValue = append(recordValue, byte(0))
	recordValue = append(recordValue, schemaIDBytes...)
	recordValue = append(recordValue, valueBytes...)
	return recordValue
}

const (
	UUID                    = "uuid"
	OWNER_USER_UUID         = "owner_user_uuid"
	USER_UUID               = "user_uuid"
	ORGANIZATION_UUID       = "organization_uuid"
	SCHOOL_UUID             = "school_uuid"
	CLASS_UUID              = "class_uuid"
	ORGANIZATION_ROLE_UUIDS = "organization_role_uuids"
	NAME                    = "name"
	PROGRAM_UUIDS           = "program_uuids"
	GIVEN_NAME              = "user_given_name"
	FAMILY_NAME             = "user_family_name"
	EMAIL                   = "user_email"
	PHONE_NUMBER            = "user_phone_number"
	DATE_OF_BIRTH           = "user_date_of_birth"
	GENDER                  = "user_gender"
	SUBJECT_UUIDS           = "subject_uuids"
	GRADE_UUIDS             = "grade_uuids"
	AGE_RANGE_UUIDS         = "age_range_uuids"
	ACADEMIC_TERM_UUID      = "academic_term_uuid"
	PARTICIPATING_AS        = "participating_as"
)

var (
	OrganizationHeaders = []string{UUID, NAME, OWNER_USER_UUID}
	SchoolHeaders       = []string{UUID, ORGANIZATION_UUID, NAME, PROGRAM_UUIDS}
	ClassHeaders        = []string{UUID, ORGANIZATION_UUID, NAME}
	OrgMemHeaders       = []string{ORGANIZATION_UUID, USER_UUID, ORGANIZATION_ROLE_UUIDS}
	SchoolMemHeaders    = []string{SCHOOL_UUID, USER_UUID}
	ClassRosterHeaders  = []string{CLASS_UUID, USER_UUID, PARTICIPATING_AS}
	UserHeaders         = []string{UUID, GIVEN_NAME, FAMILY_NAME, EMAIL, PHONE_NUMBER, DATE_OF_BIRTH, GENDER}
	ClassDetailsHeaders = []string{CLASS_UUID, SCHOOL_UUID, PROGRAM_UUIDS, SUBJECT_UUIDS, GRADE_UUIDS, AGE_RANGE_UUIDS, ACADEMIC_TERM_UUID}
)

func GetOrganizationSchemaId(schemaRegistryClient *SchemaRegistry, organizationTopic string) int {
	schemaBody := avrogen.Organization.Schema(avrogen.NewOrganization())
	return schemaRegistryClient.GetSchemaId(schemaBody, organizationTopic)
}

func GetSchoolSchemaId(schemaRegistryClient *SchemaRegistry, schoolTopic string) int {
	schemaBody := avrogen.School.Schema(avrogen.NewSchool())
	return schemaRegistryClient.GetSchemaId(schemaBody, schoolTopic)
}

func GetUserSchemaId(schemaRegistryClient *SchemaRegistry, userTopic string) int {
	schemaBody := avrogen.User.Schema(avrogen.NewUser())
	return schemaRegistryClient.GetSchemaId(schemaBody, userTopic)
}

func GetClassSchemaId(schemaRegistryClient *SchemaRegistry, classTopic string) int {
	schemaBody := avrogen.Class.Schema(avrogen.NewClass())
	return schemaRegistryClient.GetSchemaId(schemaBody, classTopic)
}

func GetClassRosterSchemaId(schemaRegistryClient *SchemaRegistry, classRosterTopic string) int {
	schemaBody := avrogen.Class.Schema(avrogen.NewClass())
	return schemaRegistryClient.GetSchemaId(schemaBody, classRosterTopic)
}

func GetOrgMemSchemaId(schemaRegistryClient *SchemaRegistry, orgMemTopic string) int {
	schemaBody := avrogen.OrganizationMembership.Schema(avrogen.NewOrganizationMembership())
	return schemaRegistryClient.GetSchemaId(schemaBody, orgMemTopic)
}

func GetClassDetailsSchemaId(schemaRegistryClient *SchemaRegistry, classDetailsTopic string) int {
	schemaBody := avrogen.ClassDetails.Schema(avrogen.NewClassDetails())
	return schemaRegistryClient.GetSchemaId(schemaBody, classDetailsTopic)
}

func GetSchoolMemSchemaId(schemaRegistryClient *SchemaRegistry, schoolMemTopic string) int {
	schemaBody := avrogen.SchoolMembership.Schema(avrogen.NewSchoolMembership())
	return schemaRegistryClient.GetSchemaId(schemaBody, schoolMemTopic)
}

func InitAvroOperations(schemaRegistryClient *SchemaRegistry) Operations {
	organizationTopic := instrument.MustGetEnv("ORGANIZATION_AVRO_TOPIC")
	schoolTopic := instrument.MustGetEnv("SCHOOL_AVRO_TOPIC")
	userTopic := instrument.MustGetEnv("USER_AVRO_TOPIC")
	classTopic := instrument.MustGetEnv("CLASS_AVRO_TOPIC")
	orgMemTopic := instrument.MustGetEnv("ORGANIZATION_MEMBERSHIP_AVRO_TOPIC")
	schoolMemTopic := instrument.MustGetEnv("SCHOOL_MEMBERSHIP_AVRO_TOPIC")
	classRosterTopic := instrument.MustGetEnv("CLASS_ROSTER_AVRO_TOPIC")
	classDetailsTopic := instrument.MustGetEnv("CLASS_DETAILS_AVRO_TOPIC")

	return Operations{
		OperationMap: map[string]Operation{
			"ORGANIZATION": {
				Topic:        organizationTopic,
				Key:          "",
				SchemaID:     GetOrganizationSchemaId(schemaRegistryClient, organizationTopic),
				SerializeRow: RowToOrganizationAvro,
				Headers:      OrganizationHeaders,
			},
			"SCHOOL": {
				Topic:        schoolTopic,
				Key:          "",
				SchemaID:     GetSchoolSchemaId(schemaRegistryClient, schoolTopic),
				SerializeRow: RowToSchoolAvro,
				Headers:      SchoolHeaders,
			},
			"USER": {
				Topic:        userTopic,
				Key:          "",
				SchemaID:     GetUserSchemaId(schemaRegistryClient, userTopic),
				SerializeRow: RowToUserAvro,
				Headers:      UserHeaders,
			},
			"CLASS": {
				Topic:        classTopic,
				Key:          "",
				SchemaID:     GetClassSchemaId(schemaRegistryClient, classTopic),
				SerializeRow: RowToClassAvro,
				Headers:      ClassHeaders,
			},
			"ORGANIZATION_MEMBERSHIP": {
				Topic:        orgMemTopic,
				Key:          "",
				SchemaID:     GetOrgMemSchemaId(schemaRegistryClient, orgMemTopic),
				SerializeRow: RowToOrgMemAvro,
				Headers:      OrgMemHeaders,
			},
			"CLASS_DETAILS": {
				Topic:        classDetailsTopic,
				Key:          "",
				SchemaID:     GetClassDetailsSchemaId(schemaRegistryClient, classDetailsTopic),
				SerializeRow: RowToClassDetailsAvro,
				Headers:      ClassDetailsHeaders,
			},
			"SCHOOL_MEMBERSHIP": {
				Topic:        schoolMemTopic,
				Key:          "",
				SchemaID:     GetSchoolMemSchemaId(schemaRegistryClient, schoolMemTopic),
				SerializeRow: RowToSchoolMemAvro,
				Headers:      SchoolMemHeaders,
			},
			"CLASS_ROSTER": {
				Topic:        classRosterTopic,
				Key:          "",
				SchemaID:     GetClassRosterSchemaId(schemaRegistryClient, classRosterTopic),
				SerializeRow: RowToClassRosterAvro,
				Headers:      ClassRosterHeaders,
			},
		},
	}
}

func RowToOrganizationAvro(row []string, tracking_uuid string, schemaId int, headerIndexes map[string]int) ([]byte, error) {
	// Takes a slice of columns representing an organization and encodes to avro bytes
	md := avrogen.OrganizationMetadata{
		Origin_application: os.Getenv("METADATA_ORIGIN_APPLICATION"),
		Region:             os.Getenv("METADATA_REGION"),
		Tracking_uuid:      tracking_uuid,
	}
	pl := avrogen.OrganizationPayload{
		Uuid:            row[headerIndexes[UUID]],
		Name:            row[headerIndexes[NAME]],
		Owner_user_uuid: row[headerIndexes[OWNER_USER_UUID]],
	}
	codec := avrogen.Organization{Payload: pl, Metadata: md}
	return serializeAvroRecord(codec, schemaId), nil
}

func RowToSchoolAvro(row []string, tracking_uuid string, schemaId int, headerIndexes map[string]int) ([]byte, error) {
	// Takes a slice of columns representing a school and encodes to avro bytes
	md := avrogen.SchoolMetadata{
		Origin_application: os.Getenv("METADATA_ORIGIN_APPLICATION"),
		Region:             os.Getenv("METADATA_REGION"),
		Tracking_uuid:      tracking_uuid,
	}
	pl := avrogen.SchoolPayload{
		Uuid:              makeAvroOptionalString(row[headerIndexes[UUID]]),
		Organization_uuid: row[headerIndexes[ORGANIZATION_UUID]],
		Name:              row[headerIndexes[NAME]],
		Program_uuids:     makeAvroOptionalArrayString(row[headerIndexes[PROGRAM_UUIDS]]),
	}

	codec := avrogen.School{Payload: pl, Metadata: md}
	return serializeAvroRecord(codec, schemaId), nil
}

func RowToUserAvro(row []string, tracking_uuid string, schemaId int, headerIndexes map[string]int) ([]byte, error) {
	// Takes a slice of columns representing a user and encodes to avro bytes
	md := avrogen.UserMetadata{
		Origin_application: os.Getenv("METADATA_ORIGIN_APPLICATION"),
		Region:             os.Getenv("METADATA_REGION"),
		Tracking_uuid:      tracking_uuid,
	}
	pl := avrogen.UserPayload{
		Uuid:          row[headerIndexes[UUID]],
		Given_name:    row[headerIndexes[GIVEN_NAME]],
		Family_name:   row[headerIndexes[FAMILY_NAME]],
		Gender:        row[headerIndexes[GENDER]],
		Phone_number:  makeAvroOptionalString(row[headerIndexes[PHONE_NUMBER]]),
		Email:         makeAvroOptionalString(row[headerIndexes[EMAIL]]),
		Date_of_birth: makeAvroOptionalString(row[headerIndexes[DATE_OF_BIRTH]]),
	}

	codec := avrogen.User{Payload: pl, Metadata: md}
	return serializeAvroRecord(codec, schemaId), nil
}

func RowToClassAvro(row []string, tracking_uuid string, schemaId int, headerIndexes map[string]int) ([]byte, error) {
	// Takes a slice of columns representing a class and encodes to avro bytes
	md := avrogen.ClassMetadata{
		Origin_application: os.Getenv("METADATA_ORIGIN_APPLICATION"),
		Region:             os.Getenv("METADATA_REGION"),
		Tracking_uuid:      tracking_uuid,
	}

	pl := avrogen.ClassPayload{
		Uuid:              makeAvroOptionalString(row[headerIndexes[UUID]]),
		Name:              row[headerIndexes[NAME]],
		Organization_uuid: row[headerIndexes[ORGANIZATION_UUID]],
	}

	codec := avrogen.Class{Payload: pl, Metadata: md}
	return serializeAvroRecord(codec, schemaId), nil
}

func RowToClassRosterAvro(row []string, tracking_uuid string, schemaId int, headerIndexes map[string]int) ([]byte, error) {
	// Takes a slice of columns representing a class and encodes to avro bytes
	md := avrogen.ClassRosterMetadata{
		Origin_application: os.Getenv("METADATA_ORIGIN_APPLICATION"),
		Region:             os.Getenv("METADATA_REGION"),
		Tracking_uuid:      tracking_uuid,
	}

	pl := avrogen.ClassRosterPayload{
		Class_uuid:       row[headerIndexes[CLASS_UUID]],
		User_uuid:        row[headerIndexes[USER_UUID]],
		Participating_as: row[headerIndexes[PARTICIPATING_AS]],
	}

	codec := avrogen.ClassRoster{Payload: pl, Metadata: md}
	return serializeAvroRecord(codec, schemaId), nil
}

func RowToOrgMemAvro(row []string, tracking_uuid string, schemaId int, headerIndexes map[string]int) ([]byte, error) {
	// Takes a slice of columns representing a class and encodes to avro bytes
	orgRoleUuids := strings.Split(row[headerIndexes[ORGANIZATION_ROLE_UUIDS]], ";")
	md := avrogen.OrganizationMembershipMetadata{
		Origin_application: os.Getenv("METADATA_ORIGIN_APPLICATION"),
		Region:             os.Getenv("METADATA_REGION"),
		Tracking_uuid:      tracking_uuid,
	}

	pl := avrogen.OrganizationMembershipPayload{
		Organization_uuid:       row[headerIndexes[ORGANIZATION_UUID]],
		User_uuid:               row[headerIndexes[USER_UUID]],
		Organization_role_uuids: orgRoleUuids,
	}

	codec := avrogen.OrganizationMembership{Payload: pl, Metadata: md}
	return serializeAvroRecord(codec, schemaId), nil
}

func RowToClassDetailsAvro(row []string, tracking_uuid string, schemaId int, headerIndexes map[string]int) ([]byte, error) {
	// Takes a slice of columns representing a class and encodes to avro bytes
	md := avrogen.ClassDetailsMetadata{
		Origin_application: os.Getenv("METADATA_ORIGIN_APPLICATION"),
		Region:             os.Getenv("METADATA_REGION"),
		Tracking_uuid:      tracking_uuid,
	}
	pl := avrogen.ClassDetailsPayload{
		Class_uuid:         row[headerIndexes[CLASS_UUID]],
		School_uuid:        makeAvroOptionalString(row[headerIndexes[SCHOOL_UUID]]),
		Program_uuids:      makeAvroOptionalArrayString(row[headerIndexes[PROGRAM_UUIDS]]),
		Subject_uuids:      makeAvroOptionalArrayString(row[headerIndexes[SUBJECT_UUIDS]]),
		Grade_uuids:        makeAvroOptionalArrayString(row[headerIndexes[GRADE_UUIDS]]),
		Age_range_uuids:    makeAvroOptionalArrayString(row[headerIndexes[AGE_RANGE_UUIDS]]),
		Academic_term_uuid: makeAvroOptionalString(row[headerIndexes[ACADEMIC_TERM_UUID]]),
	}

	codec := avrogen.ClassDetails{Payload: pl, Metadata: md}
	return serializeAvroRecord(codec, schemaId), nil
}

func RowToSchoolMemAvro(row []string, tracking_uuid string, schemaId int, headerIndexes map[string]int) ([]byte, error) {
	// Takes a slice of columns representing a class and encodes to avro bytes
	md := avrogen.SchoolMembershipMetadata{
		Origin_application: os.Getenv("METADATA_ORIGIN_APPLICATION"),
		Region:             os.Getenv("METADATA_REGION"),
		Tracking_uuid:      tracking_uuid,
	}

	pl := avrogen.SchoolMembershipPayload{
		School_uuid: row[headerIndexes[SCHOOL_UUID]],
		User_uuid:   row[headerIndexes[USER_UUID]],
	}

	codec := avrogen.SchoolMembership{Payload: pl, Metadata: md}
	return serializeAvroRecord(codec, schemaId), nil
}
