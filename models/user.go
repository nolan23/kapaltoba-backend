package models

import (
	"time"

	"github.com/zebresel-com/mongodm"
)

type User struct {
	mongodm.DocumentBase `json:",inline" bson:",inline"`
	Name                 string      `bson:"name" json:"name"`
	Email                string      `bson:"email" json:"email"`
	PhoneNumber          string      `bson:"phonenumber" json:"phonenumber"`
	BirthDate            time.Time   `bson:"birthdate" json:"birthdate"`
	Password             string      `bson:"password" json:"password"`
	ImageProfile         string      `bson:"imageprofile" json:"imageprofile"`
	TripHistory          interface{} `bson:"triphistory" json:"triphistory" model:"Trip"`
}
