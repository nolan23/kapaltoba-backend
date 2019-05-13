package repository

import (
	"context"
	"log"
	"time"

	"github.com/nolan23/kapaltoba-backend/models"
	"github.com/nolan23/kapaltoba-backend/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoUserRepository struct {
	DB             *mongo.Database
	collectionName string
}

func NewMongoUserRepository(db *mongo.Database, col string) user.Repository {
	return &mongoUserRepository{db, col}
}

func (m *mongoUserRepository) fetch(ctx context.Context, query interface{}, limit int, skip int, sort interface{}) (res []*models.User, nextSkip int, err error) {
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
	var resu []*models.User
	for cur.Next(ctx) {
		tr := &models.User{}
		err = cur.Decode(tr)
		if err != nil {
			log.Println("error decode " + err.Error())
		}
		resu = append(resu, tr)
	}
	return resu, limit + skip, nil
}

func (m *mongoUserRepository) Fetch(ctx context.Context, limit int, skip int, sort string) (res []*models.User, nextSkip int, err error) {
	query := bson.M{"deleted": false}
	if sort == "" {
		sort = "modifiedAt"
	}
	bsonSort := bson.M{sort: 1}
	users, next, err := m.fetch(ctx, query, limit, skip, bsonSort)
	if err != nil {
		log.Println("fetch trans " + err.Error())
		return nil, 0, err
	}
	return users, next, nil
}

func (m *mongoUserRepository) fetchOne(ctx context.Context, query interface{}) (*models.User, error) {
	var result models.User
	err := m.DB.Collection(m.collectionName).FindOne(ctx, query).Decode(&result)
	if err != nil {
		log.Println("error find by id " + err.Error())
		return nil, err
	}
	return &result, nil
}

func (m *mongoUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	var result *models.User
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

func (m *mongoUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var result *models.User
	filter := bson.D{{"name", username}}
	result, err := m.fetchOne(ctx, filter)
	if err != nil {
		log.Println("error find by username " + err.Error())
		return nil, err
	}
	return result, nil
}
func (m *mongoUserRepository) Update(ctx context.Context, selector interface{}, update interface{}) error {
	updateResult, err := m.DB.Collection(m.collectionName).UpdateOne(ctx, selector, update)
	if err != nil {
		log.Println("error update user " + err.Error())
		return err
	}
	log.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	return nil
}

func (m *mongoUserRepository) Store(ctx context.Context, user *models.User) error {
	user.ID = primitive.NewObjectID()
	user.Deleted = false
	user.CreatedAt = time.Now()
	user.ModifiedAt = time.Now()
	insertResult, err := m.DB.Collection(m.collectionName).InsertOne(ctx, user)
	if err != nil {
		log.Println("error store transaction " + err.Error())
		return err
	}
	log.Println("Inserted  document: ", insertResult.InsertedID)
	return nil
}
func (m *mongoUserRepository) Delete(ctx context.Context, id string) error {
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
