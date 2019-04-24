package usecase

import (
	"context"
	"time"

	"github.com/nolan23/kapaltoba-backend/models"
	"github.com/nolan23/kapaltoba-backend/user"
)

type userUsecase struct {
	userRepo       user.Repository
	contextTimeout time.Duration
}

func NewUserUsecase(a user.Repository, timeout time.Duration) user.Usecase {
	return &userUsecase{
		userRepo:       a,
		contextTimeout: timeout,
	}
}

func (u *userUsecase) Fetch(ctx context.Context, limit int, skip int, sort string) ([]*models.User, int, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	listUser, nextSkip, err := u.userRepo.Fetch(ctx, limit, skip, sort)
	if err != nil {
		return nil, 0, err
	}
	return listUser, nextSkip, nil
}
func (u *userUsecase) GetByID(ctx context.Context, id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	res, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u *userUsecase) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	res, err := u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (u *userUsecase) Update(ctx context.Context, selector interface{}, update interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	err := u.userRepo.Update(ctx, selector, update)
	if err != nil {
		return err
	}
	return nil
}

func (u *userUsecase) Store(ctx context.Context, user *models.User) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	err := u.userRepo.Store(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (u *userUsecase) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	err := u.userRepo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
