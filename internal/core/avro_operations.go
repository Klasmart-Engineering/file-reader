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
)

var (
	OrganizationHeaders = []string{UUID, ORGANIZATION_NAME, OWNER_USER_ID}
	SchoolHeaders       = []string{UUID, ORGANIZATION_UUID, SCHOOL_NAME, PROGRAM_IDS}
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

func GetClassSchemaId(schemaRegistryClient *SchemaRegistry, classTopic string) int {
	schemaBody := avrogen.Class.Schema(avrogen.NewClass())
	return schemaRegistryClient.GetSchemaId(schemaBody, classTopic)
}

func InitAvroOperations(schemaRegistryClient *SchemaRegistry) Operations {
	organizationTopic := instrument.MustGetEnv("ORGANIZATION_AVRO_TOPIC")
	schoolTopic := instrument.MustGetEnv("SCHOOL_AVRO_TOPIC")
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
		Organization_id: row[headerIndexes[ORGANIZATION_UUID]],
		Name:            row[headerIndexes[SCHOOL_NAME]],
		Program_ids:     makeAvroOptionalArrayString(row[headerIndexes[PROGRAM_IDS]]),
	}

	codec := avrogen.School{Payload: pl, Metadata: md}
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
