package usecase

import (
	"context"
	"log"
	"time"

	"github.com/nolan23/kapaltoba-backend/models"
	"github.com/nolan23/kapaltoba-backend/transaction"
	"gopkg.in/mgo.v2/bson"
)

type transactionUsecase struct {
	transRepo      transaction.Repository
	contextTimeout time.Duration
}

func NewTransactionUsecase(t transaction.Repository, timeout time.Duration) transaction.Usecase {
	return &transactionUsecase{
		transRepo:      t,
		contextTimeout: timeout,
	}
}

func (ts *transactionUsecase) Fetch(ctx context.Context, limit int, skip int, sort string) (res []*models.Transaction, nextSkip int, err error) {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()
	listTransaction, nextSkip, err := ts.transRepo.Fetch(ctx, limit, skip, sort)
	if err != nil {
		log.Println("error fetch transaction usecase " + err.Error())
		return nil, 0, err
	}
	return listTransaction, nextSkip, nil
}

func (ts *transactionUsecase) FindBy(ctx context.Context, userID string, tripID string) (*models.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()
	res, err := ts.transRepo.FindBy(ctx, userID, tripID)
	if err != nil {
		log.Println("error find by transaction usecase " + err.Error())
		return nil, err
	}
	return res, nil
}
func (ts *transactionUsecase) GetByID(ctx context.Context, id string) (*models.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()

	res, err := ts.transRepo.GetByID(ctx, id)
	if err != nil {
		log.Println("error getById transaction usecase " + err.Error())
		return nil, err
	}
	return res, nil
}

func (ts *transactionUsecase) GetByUserId(ctx context.Context, userId string) (*models.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()

	res, err := ts.transRepo.GetByUserId(ctx, userId)
	if err != nil {
		log.Println("error getByUserId transaction usecase " + err.Error())
		return nil, err
	}
	return res, nil
}
func (ts *transactionUsecase) GetByTripId(ctx context.Context, tripId string) (*models.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()

	res, err := ts.transRepo.GetByTripId(ctx, tripId)
	if err != nil {
		log.Println("error getByTripId transaction usecase " + err.Error())
		return nil, err
	}
	return res, nil
}
func (ts *transactionUsecase) GetByUsername(ctx context.Context, username string) ([]*models.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()

	res, err := ts.transRepo.GetByUsername(ctx, username)
	if err != nil {
		log.Println("error getByUsername transaction usecase " + err.Error())
		return nil, err
	}
	return res, nil
}
func (ts *transactionUsecase) Update(ctx context.Context, selector interface{}, update *models.Transaction) error {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()
	update.ModifiedAt = time.Now()
	err := ts.transRepo.Update(ctx, selector, bson.M{"$set": &update})
	if err != nil {
		log.Println("error update trip usecase " + err.Error())
		return err
	}
	return nil
}

func (ts *transactionUsecase) Store(ctx context.Context, transaction *models.Transaction) error {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()

	err := ts.transRepo.Store(ctx, transaction)
	if err != nil {
		log.Println("error store transaction usecase " + err.Error())
		return err
	}
	return nil
}
func (ts *transactionUsecase) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()
	err := ts.transRepo.Delete(ctx, id)
	if err != nil {
		log.Println("error delete transaction usecase " + err.Error())
		return err
	}
	return nil
}
