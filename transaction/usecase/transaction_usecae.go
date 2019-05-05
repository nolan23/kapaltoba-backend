package usecase

import (
	"context"
	"log"
	"time"

	"github.com/nolan23/kapaltoba-backend/models"
	"github.com/nolan23/kapaltoba-backend/transaction"
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
func (ts *transactionUsecase) Update(ctx context.Context, selector interface{}, update interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, ts.contextTimeout)
	defer cancel()
	err := ts.transRepo.Update(ctx, selector, update)
	if err != nil {
		log.Println("error update transaction usecase " + err.Error())
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
