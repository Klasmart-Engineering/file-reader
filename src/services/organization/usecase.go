package organization

import (
	"context"
	"file_reader/src/models"

	"github.com/riferrei/srclient"
)

type UseCase interface {
	Create(ctx context.Context, org *models.Organization) (*models.Organization, error)
	PublishCreate(ctx context.Context, org *models.Organization, schema *srclient.Schema) error
	PublishUpdate(ctx context.Context, org *models.Organization) error
}
