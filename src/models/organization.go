package models

import (
	"file_reader/src/protos"
)

// Organization model

type Organization struct {
	Guid string
	Name string
}

// Convert Organization to proto
func (o *Organization) ToProto() *protos.Organization {
	return &protos.Organization{
		Guid: o.Guid,
		Name: o.Name,
	}
}

func OrganizationFromProto(o *protos.Organization) (*Organization, error) {
	return &Organization{
		Guid: o.GetGuid(),
		Name: o.GetName(),
	}, nil

}
