package usecase

import (
	"context"
	"log"
	"time"

	"github.com/nolan23/kapaltoba-backend/boat"

	"github.com/nolan23/kapaltoba-backend/models"
	"github.com/nolan23/kapaltoba-backend/trip"
	"github.com/nolan23/kapaltoba-backend/user"
	"gopkg.in/mgo.v2/bson"
)

type tripUsecase struct {
	tripRepo       trip.Repository
	userRepo       user.Repository
	boatRepo       boat.Repository
	contextTimeout time.Duration
}

func NewTripUsecase(t trip.Repository, us user.Repository, br boat.Repository, timeout time.Duration) trip.Usecase {
	return &tripUsecase{
		tripRepo:       t,
		userRepo:       us,
		boatRepo:       br,
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
func (ts *tripUsecase) Update(ctx context.Context, selector interface{}, update *models.Trip) error {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()
	update.ModifiedAt = time.Now()
	err := ts.tripRepo.Update(ctx, selector, bson.M{"$set": &update})
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

func (ts *tripUsecase) GetBoat(ctx context.Context, idBoat string) (boat *models.Boat, err error) {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()
	log.Println("id " + idBoat)
	boat, err = ts.boatRepo.GetByID(ctx, idBoat)
	if err != nil {
		log.Println("error getboat trip usecase " + err.Error())
		return nil, err
	}
	return boat, nil
}

func (ts *tripUsecase) GetPassengers(ctx context.Context, idTrip string) ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()
	trip, er := ts.GetByID(ctx, idTrip)
	if er != nil {
		log.Println("error get passenger get id trip" + er.Error())
		return nil, er
	}
	var passengers []*models.User

	for _, user := range trip.Passengers.([]string) {
		passenger, err := ts.userRepo.GetByID(ctx, user)
		if err != nil {
			log.Println("error get user in trip usecase " + err.Error())
		}
		passengers = append(passengers, passenger)
	}

	return passengers, nil
}

func (ts *tripUsecase) AddPassenger(ctx context.Context, selector interface{}, trip *models.Trip, passengerId string) (*models.Trip, error) {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()
	trip.Passengers = append(trip.Passengers.([]string), passengerId)
	if trip.Available == 0 {
		return nil, nil
	}
	trip.Available = trip.Available - 1
	trip.Purchased = trip.Purchased + 1
	err := ts.Update(ctx, selector, trip)
	if err != nil {
		log.Println("error in update " + err.Error())
		return nil, err
	}
	return trip, nil
}
