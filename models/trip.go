package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Trip struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Boat        interface{}        `json:"boat" bson:"boat"`
	Origin      string             `json:"origin" bson:"origin"`
	Destination string             `json:"destination" bson:"destination"`
	Time        time.Time          `json:"time" bson:"time"`
	Status      string             `json:"status" bson:"status"`
	Price       float64            `json:"price" bson:"price"`
	Duration    time.Duration      `json:"duration" bson:"duration"`
	Available   int                `json:"available" bson:"available"`
	Purchased   int                `json:"purchased" bson:"purchased"`
	Passengers  interface{}        `json:"passengers" bson:"passengers" relation:"1n"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	ModifiedAt  time.Time          `json:"modifiedAt" bson:"modifiedAt"`
	Deleted     bool               `json:"deleted" bson:"deleted"`
}
