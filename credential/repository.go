package credential

import (
	"context"

	"github.com/nolan23/kapaltoba-backend/models"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*models.Credential, error)
	GetByUsername(ctx context.Context, username string) (*models.Credential, error)
	Update(ctx context.Context, selector interface{}, update interface{}) error
	Store(ctx context.Context, credential *models.Credential) error
	Delete(ctx context.Context, id string) error
}
