package credential

import (
	"context"

	"github.com/nolan23/kapaltoba-backend/models"
)

type Usecase interface {
	GetByID(ctx context.Context, id string) (*models.Credential, error)
	GetByUsername(ctx context.Context, username string) (*models.Credential, error)
	Update(ctx context.Context, selector interface{}, update *models.Credential) error
	Store(ctx context.Context, transaction *models.Credential) (string, error)
	Delete(ctx context.Context, id string) error
}
