package filereader

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"

	avro "github.com/KL-Engineering/file-reader/avro_gencode"
	"github.com/KL-Engineering/file-reader/src"
	"github.com/KL-Engineering/file-reader/src/instrument"
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

func GetOrganizationSchemaId(schemaRegistryClient *src.SchemaRegistry, organizationTopic string) int {
	schemaBody := avro.Organization.Schema(avro.NewOrganization())
	return schemaRegistryClient.GetSchemaId(schemaBody, organizationTopic)
}

func InitAvroOperations(schemaRegistryClient *src.SchemaRegistry) Operations {
	organizationTopic := instrument.MustGetEnv("ORGANIZATION_AVRO_TOPIC")
	return Operations{
		OperationMap: map[string]Operation{
			"ORGANIZATION": {
				Topic:        organizationTopic,
				Key:          "",
				SchemaID:     GetOrganizationSchemaId(schemaRegistryClient, organizationTopic),
				SerializeRow: RowToOrganizationAvro,
			},
		},
	}
}

// ToDo: add logic for stripping header and figuring out column order
func RowToOrganizationAvro(row []string, tracking_id string, schemaId int) ([]byte, error) {
	// Takes a slice of columns representing an organization and encodes to avro bytes
	md := avro.OrganizationMetadata{
		Origin_application: os.Getenv("METADATA_ORIGIN_APPLICATION"),
		Region:             os.Getenv("METADATA_REGION"),
		Tracking_id:        tracking_id,
	}
	pl := avro.OrganizationPayload{
		Guid:              row[0],
		Organization_name: row[1],
	}
	codec := avro.Organization{Payload: pl, Metadata: md}
	return serializeAvroRecord(codec, schemaId), nil

}
