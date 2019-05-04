package transaction

import (
	"context"

	"github.com/nolan23/kapaltoba-backend/models"
)

type Repository interface {
	Fetch(ctx context.Context, limit int, skip int, sort string) (res []*models.Transaction, nextSkip int, err error)
	GetByID(ctx context.Context, id string) (*models.Transaction, error)
	GetByUsername(ctx context.Context, username string) ([]*models.Transaction, error)
	Update(ctx context.Context, selector interface{}, update interface{}) error
	Store(ctx context.Context, transaction *models.Transaction) error
	Delete(ctx context.Context, id string) error
}
