package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Boat struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	BoatName     string             `json:"boatname" bson:"boatname"`
	Captain      string             `bson:"captain" json:"captain"`
	ViceCaptains []string           `bson:"vicecaptains" json:"vicecaptains"`
	Pictures     []string           `bson:"pictures" json:"pictures"`
	Capacity     int                `bson:"capacity" json:"capacity"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
	ModifiedAt   time.Time          `json:"modifiedAt" bson:"modifiedAt"`
	Deleted      bool               `json:"deleted" bson:"deleted"`
}
