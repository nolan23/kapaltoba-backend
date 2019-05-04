package repository

import (
	"context"
	"log"

	"github.com/nolan23/kapaltoba-backend/models"
	"github.com/nolan23/kapaltoba-backend/transaction"
	"github.com/zebresel-com/mongodm"
	"gopkg.in/mgo.v2/bson"
)

type mongoDBTransactionRepository struct {
	Conn *mongodm.Connection
}

func NewMongoDBTransactionRepository(Conn *mongodm.Connection) transaction.Repository {
	return &mongoDBTransactionRepository{Conn}
}

func (m *mongoDBTransactionRepository) fetch(ctx context.Context, query interface{}, limit int, skip int, sort string) (res []*models.Transaction, nextSkip int, err error) {
	Transaction := m.Conn.Model("Transaction")
	transaction := []*models.Transaction{}
	err = Transaction.Find(query).Sort(sort).Skip(skip).Limit(limit).Exec(&transaction)
	if err != nil {
		log.Fatal("error in find transaction " + err.Error())
		return nil, skip + limit, err
	}
	return transaction, skip + limit, nil
}

func (m *mongoDBTransactionRepository) Fetch(ctx context.Context, limit int, skip int, sort string) (res []*models.Transaction, nextSkip int, err error) {
	query := bson.M{"deleted": false}
	if sort == "" {
		sort = "_modifiedAt"
	}
	result, nextSkip, err := m.fetch(ctx, query, limit, skip, sort)
	if err != nil {
		log.Fatal("error in fetch transaction " + err.Error())
		return nil, nextSkip, err
	}
	return result, nextSkip, nil
}

func (m *mongoDBTransactionRepository) GetByID(ctx context.Context, id string) (*models.Transaction, error) {
	Transaction := m.Conn.Model("Transaction")
	transaction := &models.Transaction{}
	err := Transaction.FindId(bson.ObjectIdHex(id)).Exec(transaction)
	if _, ok := err.(*mongodm.NotFoundError); ok {
		return nil, err
	} else if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (m *mongoDBTransactionRepository) GetByUsername(ctx context.Context, username string) ([]*models.Transaction, error) {
	Transaction := m.Conn.Model("Transaction")
	transactions := []*models.Transaction{}
	returnedTransactions := []*models.Transaction{}
	err := Transaction.Find(bson.M{"deleted": false}).Populate("User").Exec(transactions)
	if err != nil {
		log.Fatal("error in get by username " + err.Error())
		return nil, err
	}
	for _, transaction := range transactions {
		if user, ok := transaction.User.(*models.User); ok {
			if user.Name == username {
				returnedTransactions = append(returnedTransactions, transaction)
			}
		}
	}
	return returnedTransactions, nil
}

func (m *mongoDBTransactionRepository) Update(ctx context.Context, selector interface{}, update interface{}) error {
	Transaction := m.Conn.Model("Transaction")
	err := Transaction.Update(selector, update)
	if err != nil {
		log.Fatal("error update transaction " + err.Error())
		return err
	}
	return nil
}

func (m *mongoDBTransactionRepository) Store(ctx context.Context, transaction *models.Transaction) error {
	Transaction := m.Conn.Model("Transaction")
	Transaction.New(transaction)
	err := transaction.Save()
	if err != nil {
		log.Fatal("error in stre transaction " + err.Error())
		return err
	}
	return nil
}

func (m *mongoDBTransactionRepository) Delete(ctx context.Context, id string) error {
	transaction, err := m.GetByID(ctx, id)
	if err != nil {
		return err
	}
	err = transaction.Delete()
	if err != nil {
		log.Fatal("error in delete transaction ")
		return err
	}
	return nil
}
