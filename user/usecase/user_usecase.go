package usecase

import (
	"context"
	"log"
	"time"

	"github.com/nolan23/kapaltoba-backend/credential"
	"github.com/nolan23/kapaltoba-backend/transaction"
	"github.com/nolan23/kapaltoba-backend/trip"

	"github.com/nolan23/kapaltoba-backend/models"
	"github.com/nolan23/kapaltoba-backend/user"
)

type userUsecase struct {
	userRepo       user.Repository
	tripRepo       trip.Repository
	transRepo      transaction.Repository
	credentialRepo credential.Repository
	contextTimeout time.Duration
}

func NewUserUsecase(a user.Repository, tr trip.Repository, trr transaction.Repository, cr credential.Repository, timeout time.Duration) user.Usecase {
	return &userUsecase{
		userRepo:       a,
		tripRepo:       tr,
		transRepo:      trr,
		credentialRepo: cr,
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

	cred, err := u.credentialRepo.GetByUsername(ctx, username)
	if err != nil {
		log.Println("error get cred in user usecase " + err.Error())
		return nil, err
	}
	usr, er := u.userRepo.GetByCredID(ctx, cred.ID.Hex())
	if er != nil {
		log.Println("error get user by credential " + er.Error())
		return nil, er
	}
	return usr, nil
}

func (u *userUsecase) GetUserTrips(ctx context.Context, id string) ([]*models.Trip, error) {
	usr, err := u.GetByID(ctx, id)
	if err != nil {
		log.Println("error get user in get user trips usecase " + err.Error())
		return nil, err
	}
	if usr.TripHistory == nil {
		return nil, nil
	}
	var userTrips []*models.Trip
	for _, trip := range usr.TripHistory.([]string) {
		userTrip, err := u.tripRepo.GetByID(ctx, trip)
		if err != nil {
			log.Println("error get user in trip usecase " + err.Error())
		}
		userTrips = append(userTrips, userTrip)
	}

	return userTrips, nil
}

func (u *userUsecase) GetTransactions(ctx context.Context, id string) ([]*models.Transaction, error) {
	usr, err := u.GetByID(ctx, id)
	if err != nil {
		log.Println("error get user in get user trips usecase " + err.Error())
		return nil, err
	}
	if usr.Transactions == nil {
		return nil, nil
	}
	var transactions []*models.Transaction
	for _, t := range usr.Transactions.([]string) {
		transaction, err := u.transRepo.GetByID(ctx, t)
		if err != nil {
			log.Println("error get user in trip usecase " + err.Error())
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
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
		log.Println("error store user usecase " + err.Error())
		return err
	}
	return nil
}

func (u *userUsecase) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()
	err := u.userRepo.Delete(ctx, id)
	if err != nil {
		log.Println("error delete user usecase " + err.Error())
		return err
	}
	return nil
}
