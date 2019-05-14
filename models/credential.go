package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Credential struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	Username   string             `json:"username" bson:"username"`
	Password   string             `json:"password" bson:"password"`
	Role       string             `json:"role" bson:"role"`
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
	ModifiedAt time.Time          `json:"modifiedAt" bson:"modifiedAt"`
	Deleted    bool               `json:"deleted" bson:"deleted"`
}
