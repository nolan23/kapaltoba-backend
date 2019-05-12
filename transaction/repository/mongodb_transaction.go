package repository

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/nolan23/kapaltoba-backend/models"
	"github.com/nolan23/kapaltoba-backend/transaction"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoTransactionRepository struct {
	DB             *mongo.Database
	collectionName string
}

func NewMongoTransactionRepository(db *mongo.Database, col string) transaction.Repository {
	return &mongoTransactionRepository{db, col}
}

func (m *mongoTransactionRepository) fetch(ctx context.Context, query interface{}, limit int, skip int, sort interface{}) (res []*models.Transaction, nextSkip int, err error) {
	var limit64 = int64(limit)
	var skip64 = int64(skip)
	findOptions := options.Find()
	findOptions.SetLimit(limit64)
	findOptions.SetSkip(skip64)
	findOptions.SetSort(sort)

	cur, err := m.DB.Collection(m.collectionName).Find(ctx, query, findOptions)
	if err != nil {
		log.Println("error fetch " + err.Error())
		return nil, 0, err
	}
	var resu []*models.Transaction
	for cur.Next(ctx) {
		log.Println("masuk sinis")
		tr := &models.Transaction{}
		err = cur.Decode(tr)
		if err != nil {
			log.Println("error decode " + err.Error())
		}
		resu = append(resu, tr)
	}
	return resu, limit + skip, nil
}

func (m *mongoTransactionRepository) Fetch(ctx context.Context, limit int, skip int, sort string) (res []*models.Transaction, nextSkip int, err error) {
	query := bson.M{"deleted": false}
	bsonSort := bson.M{"modifiedAt": 1}
	trans, next, err := m.fetch(ctx, query, limit, skip, bsonSort)
	if err != nil {
		log.Println("fetch trans " + err.Error())
		return nil, 0, err
	}
	return trans, next, nil
}
func (m *mongoTransactionRepository) GetByID(ctx context.Context, id string) (*models.Transaction, error) {
	var result models.Transaction
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("error convert to ObjectID " + err.Error())
	}
	filter := bson.D{{"_id", oid}}
	err = m.DB.Collection(m.collectionName).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		log.Println("error find by id " + err.Error())
		return nil, err
	}
	return &result, nil
}
func (m *mongoTransactionRepository) GetByUsername(ctx context.Context, username string) ([]*models.Transaction, error) {
	return nil, nil
}
func (m *mongoTransactionRepository) Update(ctx context.Context, selector interface{}, update interface{}) error {
	updateResult, err := m.DB.Collection(m.collectionName).UpdateOne(ctx, selector, update)
	if err != nil {
		log.Println("error update transaction " + err.Error())
		return err
	}
	log.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	return nil
}
func (m *mongoTransactionRepository) Store(ctx context.Context, transaction *models.Transaction) error {
	transaction.Deleted = false
	transaction.CreatedAt = time.Now()
	transaction.ModifiedAt = time.Now()
	insertResult, err := m.DB.Collection(m.collectionName).InsertOne(ctx, transaction)
	if err != nil {
		log.Println("error store transaction " + err.Error())
		return err
	}
	log.Println("Inserted multiple documents: ", insertResult.InsertedID)
	return nil
}

func (m *mongoTransactionRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("error convert to ObjectID " + err.Error())
	}
	filter := bson.D{{"_id", oid}}
	update := bson.D{
		{"$set", bson.D{
			{"deleted", false},
			{"modifiedAt", time.Now()},
		}},
	}

	err = m.Update(ctx, filter, update)
	if err != nil {
		log.Println("error update in delete " + err.Error())
		return nil
	}
	return nil
}
