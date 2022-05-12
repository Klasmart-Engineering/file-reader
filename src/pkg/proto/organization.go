package proto

import (
	orgPb "file_reader/src/protos/onboarding"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

const (
	organizationProtoTopic = "organization-proto"
	orgProtoSchemaFileName = "onboarding.proto"
	organizationSchemaName = "organization"
)

var validate *validator.Validate

type ValidatedOrganization struct {
	Uuid string `validate:required,uuid4`
}

var OrganizationProto = Operation{
	topic:            organizationProtoTopic,
	rowToProtoSchema: rowToOrganizationProto,
}

func rowToOrganizationProto(row []string) (*orgPb.Organization, error) {
	md := orgPb.Metadata{
		OriginApplication: &orgPb.StringValue{Value: os.Getenv("METADATA_ORIGIN_APPLICATION")},
		Region:            &orgPb.StringValue{Value: os.Getenv("METADATA_REGION")},
		TrackingId:        &orgPb.StringValue{Value: uuid.NewString()},
	}

	// Validate uuid format
	validate = validator.New()
	validatedOrg := ValidatedOrganization{Uuid: row[0]}
	err := validate.Struct(validatedOrg)
	if err != nil {
		return nil, err
	}
	pl := orgPb.OrganizationPayload{
		Uuid: &orgPb.StringValue{Value: row[0]},
		Name: &orgPb.StringValue{Value: row[1]},
	}

	return &orgPb.Organization{Payload: &pl, Metadata: &md}, nil
}
