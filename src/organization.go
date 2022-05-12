package src

import (
	avro "file_reader/avro_gencode"
	"os"
)

const (
	OrganizationTopic = "organization"
)

func GetOrganizationSchemaIdBytes(schemaRegistryClient *SchemaRegistry) []byte {
	schemaBody := avro.Organization.Schema(avro.NewOrganization())
	return schemaRegistryClient.GetSchemaIdBytes(schemaBody, OrganizationTopic)
}

// ToDo: add logic for stripping header and figuring out column order
func RowToOrganization(row []string, tracking_id string) avroCodec {
	md := avro.OrganizationMetadata{
		Origin_application: os.Getenv("METADATA_ORIGIN_APPLICATION"),
		Region:             os.Getenv("METADATA_REGION"),
		Tracking_id:        tracking_id,
	}
	pl := avro.OrganizationPayload{
		Guid:              row[0],
		Organization_name: row[1],
	}
	return avro.Organization{Payload: pl, Metadata: md}
}
