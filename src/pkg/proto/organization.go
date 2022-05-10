package proto

import (
	orgPb "file_reader/src/protos/onboarding"
	"os"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	organizationProtoTopic = "organization_proto"
	orgProtoSchemaFileName = "onboarding.proto"
)

var OrganizationProto = Operation{
	topic:            organizationProtoTopic,
	key:              "",
	schema:           schemaRegistryClient.getProtoSchema(orgProtoSchemaFileName, organizationProtoTopic),
	rowToProtoSchema: rowToOrganizationProto,
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
