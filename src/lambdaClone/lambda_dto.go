package lambdaClone

import (
	"time"
)

type Lambda struct {
	Id         string    `bson:"_id,omitempty" json:"id"`
	Title      string    `bson:"title,omitempty" json:"title"`
	Runtime    string    `bson:"runtime,omitempty" json:"runtime"`
	Version    string    `bson:"version,omitempty" json:"version"`
	Code       string    `bson:"code,omitempty" json:"code"`
	Disabled   bool      `bson:"disabled" json:"disabled"`
	RegDate    time.Time `bson:"reg_date,omitempty" json:"reg_date"`
	UpdateDate time.Time `bson:"update_date,omitempty" json:"update_date"`
}
