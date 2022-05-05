package src

import (
	avro "file_reader/avro_gencode"
	"os"

	"github.com/google/uuid"
)

const (
	organizationTopic = "organization"
	orgSchemaFilename = "organization.avsc"
)

var Organization = Operation{
	topic:         organizationTopic,
	key:           "",
	schemaIDBytes: schemaRegistryClient.getSchemaIdBytes(orgSchemaFilename, organizationTopic),
	rowToSchema:   rowToOrganization,
}

// ToDo: add logic for stripping header and figuring out column order
func rowToOrganization(row []string) avroCodec {
	md := avro.OrganizationMetadata{
		Origin_application: os.Getenv("METADATA_ORIGIN_APPLICATION"),
		Region:             os.Getenv("METADATA_REGION"),
		Tracking_id:        uuid.NewString(),
	}
	pl := avro.OrganizationPayload{
		Guid:              row[0],
		Organization_name: row[1],
	}
	return avro.Organization{Payload: pl, Metadata: md}
}
