package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Captain struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string             `bson:"name" json"name"`
	Gender      string             `bson:"gender" json:"gender"`
	PhoneNumber string             `bson:"phone_number" json:"phoneNumber"`
	Email       string             `bson:"email" json:"email"`
	Adress      string             `bson:"adress" json:"adress"`
	Status      string             `bson:"status" json:"status"`
	Boats       []string           `bson:"boats" json:"boats"`
	Credential  string             `bson:"credential" json:"credential"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	ModifiedAt  time.Time          `json:"modifiedAt" bson:"modifiedAt"`
	Deleted     bool               `json:"deleted" bson:"deleted"`
}
