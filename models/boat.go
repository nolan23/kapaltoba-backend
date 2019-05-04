package models

import (
	"github.com/zebresel-com/mongodm"
)

type Boat struct {
	mongodm.DocumentBase `json:",inline" bson:",inline"`
	Captain              string   `bson:"captain" json:"captain"`
	ViceCaptains         []string `bson:"vicecaptains" json:"vicecaptains"`
	Pictures             []string `bson:"pictures" json:"pictures"`
	Capacity             int      `bson:"capacity" json:"capacity"`
}
