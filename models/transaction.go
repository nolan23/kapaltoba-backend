package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	User       interface{}        `json:"user" bson:"user" model:"User"`
	Trip       interface{}        `json:"trip" bson:"trip" model:"Trip"`
	Status     string             `json:"status" bson:"status"`
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
	ModifiedAt time.Time          `json:"modifiedAt" bson:"modifiedAt"`
	Deleted    bool               `json:"deleted" bson:"deleted"`
}
