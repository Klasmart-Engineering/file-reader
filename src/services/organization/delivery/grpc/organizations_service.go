package grpc

import (
	"context"
	"file_reader/src/models"
	"file_reader/src/protos"
	"file_reader/src/services/organization"

	"github.com/go-playground/validator"
	"go.uber.org/zap"
)

// organizationService grpc service
type organizationService struct {
	log            zap.Logger
	organizationUC organization.UseCase
	validate       *validator.Validate
	protos.UnimplementedOrganizationsServiceServer
}

// NewOrganizationService organizationServer constructor
func NewOrganizationService(log zap.Logger, organizationUC organization.UseCase, validate *validator.Validate) *organizationService {
	return &organizationService{log: log, organizationUC: organizationUC, validate: validate}
}

// Create new organization
func (o *organizationService) Create(ctx context.Context, req *protos.CreateRequest) (*protos.CreateResponse, error) {
	org := &models.Organization{
		Guid: req.GetOrganization().GetGuid(),
		Name: req.GetOrganization().GetName(),
	}

	created, err := o.organizationUC.Create(ctx, org)
	if err != nil {
		o.log.Error("organizationUC.Create: %v", zap.String("error", err.Error()))
		return nil, err
	}
	return &protos.CreateResponse{Organization: created.ToProto()}, nil

}
