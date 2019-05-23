package boat

import (
	"context"

	"github.com/nolan23/kapaltoba-backend/models"
)

type Usecase interface {
	Fetch(ctx context.Context, limit int, skip int, sort string) (res []*models.Boat, nextSkip int, err error)
	GetByID(ctx context.Context, id string) (*models.Boat, error)
	GetCaptain(ctx context.Context, idCaptain string) (*models.Captain, error)
	Update(ctx context.Context, selector interface{}, update *models.Boat) error
	Store(ctx context.Context, boat *models.Boat) error
	Delete(ctx context.Context, id string) error
}
