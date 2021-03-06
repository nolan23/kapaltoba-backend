package trip

import (
	"context"

	"github.com/nolan23/kapaltoba-backend/models"
)

type Repository interface {
	Fetch(ctx context.Context, limit int, skip int, sort string) (res []*models.Trip, nextSkip int, err error)
	GetByID(ctx context.Context, id string) (*models.Trip, error)
	Update(ctx context.Context, selector interface{}, update interface{}) error
	Store(ctx context.Context, trip *models.Trip) error
	Delete(ctx context.Context, id string) error
	AddPassengers(ctx context.Context, trip *models.Trip, passengers []*models.User) (*models.Trip, error)
	AddBoat(ctx context.Context, trip *models.Trip, boat *models.Boat) (*models.Trip, error)
}
