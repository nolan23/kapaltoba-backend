package repository

import (
	"context"
	"log"
	"time"

	"github.com/nolan23/kapaltoba-backend/models"
	"github.com/nolan23/kapaltoba-backend/trip"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoTripRepository struct {
	DB             *mongo.Database
	collectionName string
}

func NewMongoTripRepository(db *mongo.Database, collName string) trip.Repository {
	return &mongoTripRepository{db, collName}
}

func (m *mongoTripRepository) fetch(ctx context.Context, query interface{}, limit int, skip int, sort interface{}) (res []*models.Trip, nextSkip int, err error) {
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
	var resu []*models.Trip
	for cur.Next(ctx) {
		tr := &models.Trip{}
		err = cur.Decode(tr)
		if err != nil {
			log.Println("error decode " + err.Error())
		}
		resu = append(resu, tr)
	}
	return resu, limit + skip, nil
}

func (m *mongoTripRepository) Fetch(ctx context.Context, limit int, skip int, sort string) (res []*models.Trip, nextSkip int, err error) {
	query := bson.M{"deleted": false}
	if sort == "" {
		sort = "modifiedAt"
	}
	bsonSort := bson.M{sort: 1}
	trip, next, err := m.fetch(ctx, query, limit, skip, bsonSort)
	if err != nil {
		log.Println("fetch trans " + err.Error())
		return nil, 0, err
	}
	return trip, next, nil
}

func (m *mongoTripRepository) fetchOne(ctx context.Context, query interface{}) (*models.Trip, error) {
	var result models.Trip
	err := m.DB.Collection(m.collectionName).FindOne(ctx, query).Decode(&result)
	if err != nil {
		log.Println("error find by id " + err.Error())
		return nil, err
	}
	return &result, nil
}

func (m *mongoTripRepository) GetByID(ctx context.Context, id string) (*models.Trip, error) {
	var result *models.Trip
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("error convert to ObjectID " + err.Error())
	}
	filter := bson.D{{"_id", oid}}
	result, err = m.fetchOne(ctx, filter)
	if err != nil {
		log.Println("error find by id " + err.Error())
		return nil, err
	}
	return result, nil
}
func (m *mongoTripRepository) Update(ctx context.Context, selector interface{}, update interface{}) error {

	updateResult, err := m.DB.Collection(m.collectionName).UpdateOne(ctx, selector, update)
	if err != nil {
		log.Println("error update transaction " + err.Error())
		return err
	}
	log.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	updateResult, err = m.DB.Collection(m.collectionName).UpdateOne(ctx, selector, bson.M{"$set": bson.M{"modifiedAt": time.Now()}})
	if err != nil {
		log.Println("error update transaction " + err.Error())
		return err
	}
	log.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	return nil
}
func (m *mongoTripRepository) Store(ctx context.Context, trip *models.Trip) error {
	trip.ID = primitive.NewObjectID()
	trip.Deleted = false
	trip.CreatedAt = time.Now()
	trip.ModifiedAt = time.Now()
	insertResult, err := m.DB.Collection(m.collectionName).InsertOne(ctx, trip)
	if err != nil {
		log.Println("error store trip " + err.Error())
		return err
	}
	log.Println("Inserted  document: ", insertResult.InsertedID)
	return nil
}

func (m *mongoTripRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("error convert to ObjectID " + err.Error())
	}
	filter := bson.D{{"_id", oid}}
	update := bson.M{
		"$set": bson.M{
			"deleted":    true,
			"modifiedAt": time.Now()},
	}

	err = m.Update(ctx, filter, update)
	if err != nil {
		log.Println("error update in delete " + err.Error())
		return nil
	}
	return nil
}
func (m *mongoTripRepository) AddPassengers(ctx context.Context, trip *models.Trip, passengers []*models.User) (*models.Trip, error) {
	return nil, nil
}
func (m *mongoTripRepository) AddBoat(ctx context.Context, trip *models.Trip, boat *models.Boat) (*models.Trip, error) {
	return nil, nil
}
