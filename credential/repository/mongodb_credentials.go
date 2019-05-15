package repository

import (
	"context"
	"log"
	"time"

	"github.com/nolan23/kapaltoba-backend/credential"
	"github.com/nolan23/kapaltoba-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoCredentialRepository struct {
	DB             *mongo.Database
	collectionName string
}

func NewMongoCredentialRepository(db *mongo.Database, col string) credential.Repository {
	return &mongoCredentialRepository{db, col}
}

func (m *mongoCredentialRepository) fetchOne(ctx context.Context, query interface{}) (*models.Credential, error) {
	var result models.Credential
	err := m.DB.Collection(m.collectionName).FindOne(ctx, query).Decode(&result)
	if err != nil {
		log.Println("error find by id " + err.Error())
		return nil, err
	}
	return &result, nil
}

func (m *mongoCredentialRepository) GetByID(ctx context.Context, id string) (*models.Credential, error) {
	var result *models.Credential
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
func (m *mongoCredentialRepository) GetByUsername(ctx context.Context, username string) (*models.Credential, error) {
	var result *models.Credential
	filter := bson.D{{"username", username}}
	result, err := m.fetchOne(ctx, filter)
	if err != nil {
		log.Println("error find by username " + err.Error())
		return nil, err
	}
	return result, nil
}
func (m *mongoCredentialRepository) Update(ctx context.Context, selector interface{}, update interface{}) error {
	updateResult, err := m.DB.Collection(m.collectionName).UpdateOne(ctx, selector, update)
	if err != nil {
		log.Println("error update credential " + err.Error())
		return err
	}
	log.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	return nil
}
func (m *mongoCredentialRepository) Store(ctx context.Context, credential *models.Credential) (string, error) {
	credential.ID = primitive.NewObjectID()
	credential.Deleted = false
	credential.CreatedAt = time.Now()
	credential.ModifiedAt = time.Now()
	insertResult, err := m.DB.Collection(m.collectionName).InsertOne(ctx, credential)
	if err != nil {
		log.Println("error store credential " + err.Error())
		return "", err
	}

	log.Println("Inserted credential document: ", insertResult.InsertedID.(primitive.ObjectID).Hex())

	return insertResult.InsertedID.(primitive.ObjectID).Hex(), nil
}
func (m *mongoCredentialRepository) Delete(ctx context.Context, id string) error {
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
