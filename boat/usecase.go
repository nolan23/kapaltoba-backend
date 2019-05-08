package boat

import (
	"context"

	"github.com/nolan23/kapaltoba-backend/models"
)

type Usecase interface {
	Fetch(ctx context.Context, limit int, skip int, sort string) (res []*models.Boat, nextSkip int, err error)
	GetByID(ctx context.Context, id string) (*models.Boat, error)
	Update(ctx context.Context, selector interface{}, update interface{}) error
	Store(ctx context.Context, transaction *models.Boat) error
	Delete(ctx context.Context, id string) error
}
