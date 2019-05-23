package transaction

import (
	"context"

	"github.com/nolan23/kapaltoba-backend/models"
)

type Usecase interface {
	Fetch(ctx context.Context, limit int, skip int, sort string) (res []*models.Transaction, nextSkip int, err error)
	FindBy(ctx context.Context, userID string, tripID string) (*models.Transaction, error)
	GetByID(ctx context.Context, id string) (*models.Transaction, error)
	GetByUserId(ctx context.Context, userId string) (*models.Transaction, error)
	GetByTripId(ctx context.Context, tripId string) (*models.Transaction, error)
	GetByUsername(ctx context.Context, username string) ([]*models.Transaction, error)
	Update(ctx context.Context, selector interface{}, update interface{}) error
	Store(ctx context.Context, transaction *models.Transaction) error
	Delete(ctx context.Context, id string) error
}
