package user

import (
	"context"

	"github.com/nolan23/kapaltoba-backend/models"
)

type Repository interface {
	Fetch(ctx context.Context, limit int, skip int, sort string) (res []*models.User, nextSkip int, err error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Update(ctx context.Context, selector interface{}, update interface{}) error
	Store(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
}
