package repository

import (
	"context"
	"log"
	"time"

	"github.com/nolan23/kapaltoba-backend/boat"
	"github.com/nolan23/kapaltoba-backend/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type mongoBoatRepository struct {
	DB             *mongo.Database
	collectionName string
}

func NewMongoBoatRepository(db *mongo.Database, col string) boat.Repository {
	return &mongoBoatRepository{db, col}
}

func (m *mongoBoatRepository) fetch(ctx context.Context, query interface{}, limit int, skip int, sort interface{}) (res []*models.Boat, nextSkip int, err error) {
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
	var resu []*models.Boat
	for cur.Next(ctx) {
		tr := &models.Boat{}
		err = cur.Decode(tr)
		if err != nil {
			log.Println("error decode " + err.Error())
		}
		resu = append(resu, tr)
	}
	return resu, limit + skip, nil
}

func (m *mongoBoatRepository) Fetch(ctx context.Context, limit int, skip int, sort string) (res []*models.Boat, nextSkip int, err error) {
	query := bson.M{"deleted": false}
	if sort == "" {
		sort = "modifiedAt"
	}
	bsonSort := bson.M{sort: 1}
	boats, next, err := m.fetch(ctx, query, limit, skip, bsonSort)
	if err != nil {
		log.Println("fetch trans " + err.Error())
		return nil, 0, err
	}
	return boats, next, nil
}

func (m *mongoBoatRepository) fetchOne(ctx context.Context, query interface{}) (*models.Boat, error) {
	var result models.Boat
	err := m.DB.Collection(m.collectionName).FindOne(ctx, query).Decode(&result)
	if err != nil {
		log.Println("error find by id " + err.Error())
		return nil, err
	}
	return &result, nil
}
func (m *mongoBoatRepository) GetByID(ctx context.Context, id string) (*models.Boat, error) {
	var result *models.Boat
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("error convert to ObjectID " + err.Error())
	}
	filter := bson.M{"_id": oid}
	result, err = m.fetchOne(ctx, filter)
	if err != nil {
		log.Println("error find by id " + err.Error())
		return nil, err
	}
	return result, nil
}

func (m *mongoBoatRepository) Update(ctx context.Context, selector interface{}, update interface{}) error {
	updateResult, err := m.DB.Collection(m.collectionName).UpdateOne(ctx, selector, update)
	if err != nil {
		log.Println("error update user " + err.Error())
		return err
	}
	log.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	return nil
}
func (m *mongoBoatRepository) Store(ctx context.Context, boat *models.Boat) error {
	boat.ID = primitive.NewObjectID()
	boat.Deleted = false
	boat.CreatedAt = time.Now()
	boat.ModifiedAt = time.Now()
	insertResult, err := m.DB.Collection(m.collectionName).InsertOne(ctx, boat)
	if err != nil {
		log.Println("error store transaction " + err.Error())
		return err
	}
	log.Println("Inserted  document: ", insertResult.InsertedID)
	return nil
}
func (m *mongoBoatRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("error convert to ObjectID " + err.Error())
	}
	filter := bson.D{{"_id", oid}}
	update := bson.D{
		{"$set", bson.D{
			{"deleted", true},
			{"modifiedAt", time.Now()},
		}},
	}

	err = m.Update(ctx, filter, update)
	if err != nil {
		log.Println("error update in delete user " + err.Error())
		return nil
	}
	return nil
}
