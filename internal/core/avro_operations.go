package core

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"

	avrogen "github.com/KL-Engineering/file-reader/api/avro/avro_gencode"
	"github.com/KL-Engineering/file-reader/internal/instrument"
)

type avroCodec interface {
	// Represents the structs in the avro-gencode which have Serialize() functions
	Serialize(io.Writer) error
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
)

var (
	OrganizationHeaders = []string{UUID, ORGANIZATION_NAME, OWNER_USER_ID}
)

func GetOrganizationSchemaId(schemaRegistryClient *SchemaRegistry, organizationTopic string) int {
	schemaBody := avrogen.Organization.Schema(avrogen.NewOrganization())
	return schemaRegistryClient.GetSchemaId(schemaBody, organizationTopic)
}

func InitAvroOperations(schemaRegistryClient *SchemaRegistry) Operations {
	organizationTopic := instrument.MustGetEnv("ORGANIZATION_AVRO_TOPIC")
	return Operations{
		OperationMap: map[string]Operation{
			"ORGANIZATION": {
				Topic:        organizationTopic,
				Key:          "",
				SchemaID:     GetOrganizationSchemaId(schemaRegistryClient, organizationTopic),
				SerializeRow: RowToOrganizationAvro,
				Headers:      OrganizationHeaders,
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
