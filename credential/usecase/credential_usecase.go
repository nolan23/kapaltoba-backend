package usecase

import (
	"context"
	"log"
	"time"

	"github.com/nolan23/kapaltoba-backend/credential"
	"github.com/nolan23/kapaltoba-backend/models"
	"gopkg.in/mgo.v2/bson"
)

type credentialUsecase struct {
	credentialRepo credential.Repository
	contextTimeout time.Duration
}

func NewCredentialUsecase(t credential.Repository, timeout time.Duration) credential.Usecase {
	return &credentialUsecase{
		credentialRepo: t,
		contextTimeout: timeout,
	}
}

func (cs *credentialUsecase) GetByID(ctx context.Context, id string) (*models.Credential, error) {
	ctx, cancel := context.WithTimeout(ctx, cs.contextTimeout)
	defer cancel()
	res, err := cs.credentialRepo.GetByID(ctx, id)
	if err != nil {
		log.Println("error getById credential usecase " + err.Error())
		return nil, err
	}
	return res, nil
}
func (cs *credentialUsecase) GetByUsername(ctx context.Context, username string) (*models.Credential, error) {
	ctx, cancel := context.WithTimeout(ctx, cs.contextTimeout)
	defer cancel()

	res, err := cs.credentialRepo.GetByUsername(ctx, username)
	if err != nil {
		log.Println("error getByUsername credential usecase " + err.Error())
		return nil, err
	}
	return res, nil
}
func (cs *credentialUsecase) Update(ctx context.Context, selector interface{}, update *models.Credential) error {
	ctx, cancel := context.WithTimeout(ctx, cs.contextTimeout)
	defer cancel()
	update.ModifiedAt = time.Now()
	err := cs.credentialRepo.Update(ctx, selector, bson.M{"$set": &update})
	if err != nil {
		log.Println("error update trip usecase " + err.Error())
		return err
	}
	return nil
}
func (cs *credentialUsecase) Store(ctx context.Context, transaction *models.Credential) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, cs.contextTimeout)
	defer cancel()

	insertedId, err := cs.credentialRepo.Store(ctx, transaction)
	if err != nil {
		log.Println("error store transaction usecase " + err.Error())
		return "", err
	}
	return insertedId, nil
}
func (cs *credentialUsecase) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, cs.contextTimeout)
	defer cancel()
	err := cs.credentialRepo.Delete(ctx, id)
	if err != nil {
		log.Println("error delete transaction usecase " + err.Error())
		return err
	}
	return nil
}
