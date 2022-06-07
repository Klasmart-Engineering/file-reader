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
	UUID              = "uuid"
	ORGANIZATION_NAME = "organization_name"
	OWNER_USER_ID     = "owner_user_id"
	ORGANIZATION_UUID = "organization_id"
	SCHOOL_NAME       = "school_name"
	CLASS_NAME        = "class_name"
	PROGRAM_IDS       = "program_ids"
	GIVEN_NAME        = "user_given_name"
	FAMILY_NAME       = "user_family_name"
	EMAIL             = "user_email"
	PHONE_NUMBER      = "user_phone_number"
	DATE_OF_BIRTH     = "user_date_of_birth"
	GENDER            = "user_gender"
)

var (
	OrganizationHeaders = []string{UUID, ORGANIZATION_NAME, OWNER_USER_ID}
	SchoolHeaders       = []string{UUID, ORGANIZATION_UUID, SCHOOL_NAME, PROGRAM_IDS}
	UserHeaders         = []string{UUID, GIVEN_NAME, FAMILY_NAME, EMAIL, PHONE_NUMBER, DATE_OF_BIRTH, GENDER}
	ClassHeaders        = []string{UUID, ORGANIZATION_UUID, CLASS_NAME}
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

func InitAvroOperations(schemaRegistryClient *SchemaRegistry) Operations {
	organizationTopic := instrument.MustGetEnv("ORGANIZATION_AVRO_TOPIC")
	schoolTopic := instrument.MustGetEnv("SCHOOL_AVRO_TOPIC")
	userTopic := instrument.MustGetEnv("USER_AVRO_TOPIC")
	classTopic := instrument.MustGetEnv("CLASS_AVRO_TOPIC")
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
		},
	}
}

func RowToOrganizationAvro(row []string, tracking_id string, schemaId int, headerIndexes map[string]int) ([]byte, error) {
	// Takes a slice of columns representing an organization and encodes to avro bytes
	md := avrogen.OrganizationMetadata{
		Origin_application: os.Getenv("METADATA_ORIGIN_APPLICATION"),
		Region:             os.Getenv("METADATA_REGION"),
		Tracking_id:        tracking_id,
	}
	pl := avrogen.OrganizationPayload{
		Uuid:          row[headerIndexes[UUID]],
		Name:          row[headerIndexes[ORGANIZATION_NAME]],
		Owner_user_id: row[headerIndexes[OWNER_USER_ID]],
	}
	codec := avrogen.Organization{Payload: pl, Metadata: md}
	return serializeAvroRecord(codec, schemaId), nil
}

func RowToSchoolAvro(row []string, tracking_id string, schemaId int, headerIndexes map[string]int) ([]byte, error) {
	// Takes a slice of columns representing a school and encodes to avro bytes
	md := avrogen.SchoolMetadata{
		Origin_application: os.Getenv("METADATA_ORIGIN_APPLICATION"),
		Region:             os.Getenv("METADATA_REGION"),
		Tracking_id:        tracking_id,
	}
	pl := avrogen.SchoolPayload{
		Uuid:            makeAvroOptionalString(row[headerIndexes[UUID]]),
		Program_ids:     makeAvroOptionalArrayString(row[headerIndexes[PROGRAM_IDS]]),
		Organization_id: row[headerIndexes[ORGANIZATION_UUID]],
		Name:            row[headerIndexes[SCHOOL_NAME]],
	}

	codec := avrogen.School{Payload: pl, Metadata: md}
	return serializeAvroRecord(codec, schemaId), nil
}

func RowToUserAvro(row []string, tracking_id string, schemaId int, headerIndexes map[string]int) ([]byte, error) {
	// Takes a slice of columns representing a user and encodes to avro bytes
	md := avrogen.UserMetadata{
		Origin_application: os.Getenv("METADATA_ORIGIN_APPLICATION"),
		Region:             os.Getenv("METADATA_REGION"),
		Tracking_id:        tracking_id,
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
func RowToClassAvro(row []string, tracking_id string, schemaId int, headerIndexes map[string]int) ([]byte, error) {
	// Takes a slice of columns representing a class and encodes to avro bytes
	md := avrogen.ClassMetadata{
		Origin_application: os.Getenv("METADATA_ORIGIN_APPLICATION"),
		Region:             os.Getenv("METADATA_REGION"),
		Tracking_id:        tracking_id,
	}

	pl := avrogen.ClassPayload{
		Uuid:              makeAvroOptionalString(row[headerIndexes[UUID]]),
		Name:              row[headerIndexes[CLASS_NAME]],
		Organization_uuid: row[headerIndexes[ORGANIZATION_UUID]],
	}

	codec := avrogen.Class{Payload: pl, Metadata: md}
	return serializeAvroRecord(codec, schemaId), nil
}
