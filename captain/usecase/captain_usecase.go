package usecase

import (
	"context"
	"log"
	"time"

	"github.com/nolan23/kapaltoba-backend/trip"
	"gopkg.in/mgo.v2/bson"

	"github.com/nolan23/kapaltoba-backend/credential"

	"github.com/nolan23/kapaltoba-backend/captain"
	"github.com/nolan23/kapaltoba-backend/models"
)

type captainUsecase struct {
	captainRepo    captain.Repository
	credentialRepo credential.Repository
	tripRepo       trip.Repository
	contextTimeout time.Duration
}

func NewCaptainUsecase(t captain.Repository, cr credential.Repository, tr trip.Repository, timeout time.Duration) captain.Usecase {
	return &captainUsecase{
		captainRepo:    t,
		credentialRepo: cr,
		tripRepo:       tr,
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

func (ts *captainUsecase) GetTrips(ctx context.Context, id string) ([]*models.Trip, error) {
	captain, err := ts.GetByID(ctx, id)
	if err != nil {
		log.Println("error get captain in get trips captain usecase " + err.Error())
		return nil, err
	}
	if captain.Trips == nil {
		return nil, nil
	}
	var captainTrips []*models.Trip
	for _, trip := range captain.Trips {
		captainTrip, err := ts.tripRepo.GetByID(ctx, trip)
		if err != nil {
			log.Println("error get trip in captain usecase " + err.Error())
		}
		captainTrips = append(captainTrips, captainTrip)
	}

	return captainTrips, nil
}

func (ts *captainUsecase) Update(ctx context.Context, selector interface{}, update *models.Captain) error {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()
	update.ModifiedAt = time.Now()
	err := ts.captainRepo.Update(ctx, selector, bson.M{"$set": &update})
	if err != nil {
		log.Println("error update trip usecase " + err.Error())
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
