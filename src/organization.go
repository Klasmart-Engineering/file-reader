package src

import (
	avro "file_reader/avro_gencode"
	orgPb "file_reader/src/protos/onboarding"
	"os"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	organizationTopic      = "organization"
	organizationProtoTopic = "organization_proto"
	orgSchemaFilename      = "organization.avsc"
	orgProtoSchemaFileName = "onboarding.proto"
)

var Organization = Operation{
	topic:       organizationTopic,
	key:         "",
	rowToSchema: rowToOrganization,
}

var OrganizationProto = Operation{
	topic:            organizationProtoTopic,
	key:              "",
	schema:           schemaRegistryClient.getProtoSchema(orgProtoSchemaFileName, organizationProtoTopic),
	rowToProtoSchema: rowToOrganizationProto,
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

func rowToOrganizationProto(row []string) *orgPb.Organization {
	md := orgPb.Metadata{
		OriginApplication: wrapperspb.String(os.Getenv("METADATA_ORIGIN_APPLICATION")),
		Region:            wrapperspb.String(os.Getenv("METADATA_REGION")),
		TrackingId:        wrapperspb.String(uuid.NewString()),
	}
	pl := orgPb.OrganizationPayload{
		Uuid: wrapperspb.String(row[0]),
		Name: wrapperspb.String(row[1]),
	}
	return &orgPb.Organization{Payload: &pl, Metadata: &md}
}
