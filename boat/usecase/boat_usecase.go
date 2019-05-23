package usecase

import (
	"context"
	"log"
	"time"

	"github.com/nolan23/kapaltoba-backend/captain"

	"github.com/nolan23/kapaltoba-backend/boat"
	"github.com/nolan23/kapaltoba-backend/models"
)

type boatUsecase struct {
	boatRepo       boat.Repository
	captainRepo    captain.Repository
	contextTimeout time.Duration
}

func NewBoatUsecase(t boat.Repository, cr captain.Repository, timeout time.Duration) boat.Usecase {
	return &boatUsecase{
		boatRepo:       t,
		captainRepo:    cr,
		contextTimeout: timeout,
	}
}

func (ts *boatUsecase) Fetch(ctx context.Context, limit int, skip int, sort string) ([]*models.Boat, int, error) {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()
	listBoat, nextSkip, err := ts.boatRepo.Fetch(ctx, limit, skip, sort)
	if err != nil {
		log.Println("error fetch boat usecase " + err.Error())
		return nil, 0, err
	}
	return listBoat, nextSkip, nil
}

func (ts *boatUsecase) GetByID(ctx context.Context, id string) (*models.Boat, error) {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()

	res, err := ts.boatRepo.GetByID(ctx, id)
	if err != nil {
		log.Println("error getById boat usecase " + err.Error())
		return nil, err
	}
	return res, nil
}

func (ts *boatUsecase) GetCaptain(ctx context.Context, idCaptain string) (*models.Captain, error) {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()

	res, err := ts.captainRepo.GetByID(ctx, idCaptain)
	if err != nil {
		log.Println("error get captain boat usecase " + err.Error())
		return nil, err
	}
	return res, nil
}

func (ts *boatUsecase) Update(ctx context.Context, selector interface{}, update interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()
	err := ts.boatRepo.Update(ctx, selector, update)
	if err != nil {
		log.Println("error update boat usecase " + err.Error())
		return err
	}
	return nil
}
func (ts *boatUsecase) Store(ctx context.Context, boat *models.Boat) error {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()

	err := ts.boatRepo.Store(ctx, boat)
	if err != nil {
		log.Println("error store boat usecase " + err.Error())
		return err
	}
	return nil
}
func (ts *boatUsecase) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()
	err := ts.boatRepo.Delete(ctx, id)
	if err != nil {
		log.Println("error delete boat usecase " + err.Error())
		return err
	}
	return nil
}
