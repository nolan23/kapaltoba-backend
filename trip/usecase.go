package trip

import (
	"context"

	"github.com/nolan23/kapaltoba-backend/models"
)

type Usecase interface {
	Fetch(ctx context.Context, limit int, skip int, sort string) (res []*models.Trip, nextSkip int, err error)
	GetByID(ctx context.Context, id string) (*models.Trip, error)
	Update(ctx context.Context, selector interface{}, update *models.Trip) error
	Store(ctx context.Context, trip *models.Trip) error
	Delete(ctx context.Context, id string) error
	GetPassengers(ctx context.Context, idTrip string) (passengers []*models.User, err error)
	GetBoat(ctx context.Context, idBoat string) (boat *models.Boat, err error)
	GetCaptain(ctx context.Context, idCaptain string) (boat *models.Captain, err error)
	AddPassenger(ctx context.Context, selector interface{}, trip *models.Trip, passengerId string) (*models.Trip, error)
}
