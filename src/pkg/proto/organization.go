package proto

import (
	"file_reader/src/pkg/validation"
	orgPb "file_reader/src/protos/onboarding"
	"os"
)

const (
	orgProtoSchemaFileName = "onboarding.proto"
)

var OrganizationProto = Operation{
	rowToProtoSchema: rowToOrganizationProto,
}

func rowToOrganizationProto(row []string, trackingId string) (*orgPb.Organization, error) {
	md := orgPb.Metadata{
		OriginApplication: &orgPb.StringValue{Value: os.Getenv("METADATA_ORIGIN_APPLICATION")},
		Region:            &orgPb.StringValue{Value: os.Getenv("METADATA_REGION")},
		TrackingId:        &orgPb.StringValue{Value: trackingId},
	}

	// Validate uuid format

	validatedOrgId := validation.ValidatedOrganizationID{Uuid: row[0]}

	err := validation.UUIDValidate(validatedOrgId)
	if err != nil {
		return nil, err
	}
	pl := orgPb.OrganizationPayload{
		Uuid: &orgPb.StringValue{Value: row[0]},
		Name: &orgPb.StringValue{Value: row[1]},
	}

	return &orgPb.Organization{Payload: &pl, Metadata: &md}, nil
}
