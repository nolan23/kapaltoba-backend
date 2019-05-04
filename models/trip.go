package models

import (
	"time"

	"github.com/zebresel-com/mongodm"
)

type Trip struct {
	mongodm.DocumentBase `json:",inline" bson:",inline"`
	Boat                 interface{} `json:"boat" bson:"boat" model:"Boat"`
	Origin               string      `json:"origin" bson:"origin"`
	Destination          string      `json:"destination" bson:"destination"`
	Time                 time.Time   `json:"time" bson:"time"`
	Status               string      `json:"status" bson:"status"`
	Price                float64     `json:"price" bson:"price"`
	Duration             string      `json:"duration" bson:"duration"`
	Passenger            interface{} `json:"passenger" bson:"passenger"`
}
