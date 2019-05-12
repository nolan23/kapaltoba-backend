package models

import (
	"time"
)

type Transaction struct {
	User       interface{} `json:"user" bson:"user" model:"User"`
	Trip       interface{} `json:"trip" bson:"trip" model:"Trip"`
	Status     string      `json:"status" bson:"status"`
	CreatedAt  time.Time   `json:"createdAt" bson:"createdAt"`
	ModifiedAt time.Time   `json:"ModifiedAt" bson:"ModifiedAt"`
	Deleted    bool        `json:"deleted" bson:"deleted"`
}
