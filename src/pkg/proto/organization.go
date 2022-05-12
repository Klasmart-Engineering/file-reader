package proto

import (
	orgPb "file_reader/src/protos/onboarding"
	"os"

	"github.com/go-playground/validator/v10"
)

const (
	organizationProtoTopic = "organization-proto"
	orgProtoSchemaFileName = "onboarding.proto"
	organizationSchemaName = "organization"
)

var validate *validator.Validate

type ValidatedOrganizationID struct {
	Uuid string `validate:required,uuid4`
}
type ValidateTrackingId struct {
	Uuid string `validate:required,uuid4`
}

var OrganizationProto = Operation{
	topic:            organizationProtoTopic,
	rowToProtoSchema: rowToOrganizationProto,
}

func rowToOrganizationProto(row []string, trackingId string) (*orgPb.Organization, error) {
	md := orgPb.Metadata{
		OriginApplication: &orgPb.StringValue{Value: os.Getenv("METADATA_ORIGIN_APPLICATION")},
		Region:            &orgPb.StringValue{Value: os.Getenv("METADATA_REGION")},
		TrackingId:        &orgPb.StringValue{Value: trackingId},
	}

	// Validate uuid format
	validate = validator.New()
	validatedOrg := ValidatedOrganizationID{Uuid: row[0]}
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
