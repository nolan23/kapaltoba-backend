package models

import (
	"github.com/zebresel-com/mongodm"
)

type Transaction struct {
	mongodm.DocumentBase `json:",inline" bson:",inline"`
	User                 interface{} `json:"user" bson:"user" model:"User"`
	Trip                 interface{} `json:"trip" bson:"trip" model:"Trip"`
	Status               string      `json:"status" bson:"status"`
}
