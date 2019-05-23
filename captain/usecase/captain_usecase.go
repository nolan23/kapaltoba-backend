package usecase

import (
	"context"
	"log"
	"time"

	"github.com/nolan23/kapaltoba-backend/credential"

	"github.com/nolan23/kapaltoba-backend/captain"
	"github.com/nolan23/kapaltoba-backend/models"
)

type captainUsecase struct {
	captainRepo    captain.Repository
	credentialRepo credential.Repository
	contextTimeout time.Duration
}

func NewCaptainUsecase(t captain.Repository, cr credential.Repository, timeout time.Duration) captain.Usecase {
	return &captainUsecase{
		captainRepo:    t,
		credentialRepo: cr,
		contextTimeout: timeout,
	}
}

func (ts *captainUsecase) Fetch(ctx context.Context, limit int, skip int, sort string) ([]*models.Captain, int, error) {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()
	listCaptain, nextSkip, err := ts.captainRepo.Fetch(ctx, limit, skip, sort)
	if err != nil {
		log.Println("error fetch captain usecase " + err.Error())
		return nil, 0, err
	}
	return listCaptain, nextSkip, nil
}

func (ts *captainUsecase) GetByID(ctx context.Context, id string) (*models.Captain, error) {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()

	res, err := ts.captainRepo.GetByID(ctx, id)
	if err != nil {
		log.Println("error getById captain usecase " + err.Error())
		return nil, err
	}
	return res, nil
}

func (ts *captainUsecase) GetByUsername(ctx context.Context, username string) (*models.Captain, error) {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()

	cred, err := ts.credentialRepo.GetByUsername(ctx, username)
	if err != nil {
		log.Println("error get cred in captain usecase " + err.Error())
		return nil, err
	}
	captain, er := ts.captainRepo.GetByCredID(ctx, cred.ID.Hex())
	if er != nil {
		log.Println("error get captain by credential " + er.Error())
		return nil, er
	}
	return captain, nil
}

func (ts *captainUsecase) Update(ctx context.Context, selector interface{}, update interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()
	err := ts.captainRepo.Update(ctx, selector, update)
	if err != nil {
		log.Println("error update captain usecase " + err.Error())
		return err
	}
	return nil
}
func (ts *captainUsecase) Store(ctx context.Context, captain *models.Captain) error {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()

	err := ts.captainRepo.Store(ctx, captain)
	if err != nil {
		log.Println("error store captain usecase " + err.Error())
		return err
	}
	return nil
}
func (ts *captainUsecase) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()
	err := ts.captainRepo.Delete(ctx, id)
	if err != nil {
		log.Println("error delete captain usecase " + err.Error())
		return err
	}
	return nil
}
