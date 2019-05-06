package usecase

import (
	"context"
	"log"
	"time"

	"github.com/nolan23/kapaltoba-backend/models"
	"github.com/nolan23/kapaltoba-backend/trip"
)

type tripUsecase struct {
	tripRepo       trip.Repository
	contextTimeout time.Duration
}

func NewTripUsecase(t trip.Repository, timeout time.Duration) trip.Usecase {
	return &tripUsecase{
		tripRepo:       t,
		contextTimeout: timeout,
	}
}

func (ts *tripUsecase) Fetch(ctx context.Context, limit int, skip int, sort string) (res []*models.Trip, nextSkip int, err error) {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()
	listTrip, nextSkip, err := ts.tripRepo.Fetch(ctx, limit, skip, sort)
	if err != nil {
		log.Println("error fetch trip usecase " + err.Error())
		return nil, 0, err
	}
	return listTrip, nextSkip, nil
}
func (ts *tripUsecase) GetByID(ctx context.Context, id string) (*models.Trip, error) {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()

	res, err := ts.tripRepo.GetByID(ctx, id)
	if err != nil {
		log.Println("error getById trip usecase " + err.Error())
		return nil, err
	}
	return res, nil
}
func (ts *tripUsecase) Update(ctx context.Context, selector interface{}, update interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()
	err := ts.tripRepo.Update(ctx, selector, update)
	if err != nil {
		log.Println("error update trip usecase " + err.Error())
		return err
	}
	return nil
}
func (ts *tripUsecase) Store(ctx context.Context, trip *models.Trip) error {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()

	err := ts.tripRepo.Store(ctx, trip)
	if err != nil {
		log.Println("error store trip usecase " + err.Error())
		return err
	}
	return nil
}
func (ts *tripUsecase) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()
	err := ts.tripRepo.Delete(ctx, id)
	if err != nil {
		log.Println("error delete trip usecase " + err.Error())
		return err
	}
	return nil
}

func (ts *tripUsecase) GetPassenger(ctx context.Context, idTrip string) (passengers []*models.User, err error) {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()
	trip, er := ts.GetByID(ctx, idTrip)
	if er != nil {
		log.Println("error get passenger get id trip" + er.Error())
		return nil, er
	}
	er = trip.Populate("Passenger")
	if er != nil {
		log.Println("error in populate Passenger " + er.Error())
	}
	if users, ok := trip.Passenger.([]*models.User); ok {
		for _, user := range users {
			passengers = append(passengers, user)
		}
	} else {
		log.Println("error type assertion get passenger trip usecase ")
	}
	return passengers, nil
}
