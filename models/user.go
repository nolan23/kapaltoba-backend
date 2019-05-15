package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	Credential   interface{}        `bson:"credential" json:"credential"`
	Name         string             `bson:"name" json:"name"`
	Email        string             `bson:"email" json:"email"`
	PhoneNumber  string             `bson:"phonenumber" json:"phonenumber"`
	BirthDate    time.Time          `bson:"birthdate" json:"birthdate"`
	ImageProfile string             `bson:"imageprofile" json:"imageprofile"`
	TripHistory  interface{}        `bson:"triphistory" json:"triphistory" model:"Trip" relation:"1n"`
	Transactions interface{}        `bson:"transactions" json:"transactions"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
	ModifiedAt   time.Time          `json:"modifiedAt" bson:"modifiedAt"`
	Deleted      bool               `json:"deleted" bson:"deleted"`
}
