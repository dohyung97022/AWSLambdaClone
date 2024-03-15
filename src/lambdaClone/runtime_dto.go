package lambdaClone

import (
	"time"
)

type Runtime struct {
	Runtime     string    `bson:"runtime,omitempty" json:"runtime"`
	Version     string    `bson:"version,omitempty" json:"version"`
	Image       string    `bson:"image,omitempty" json:"image"`
	RunCommand  string    `bson:"run_command,omitempty" json:"run_command"`
	DefaultCode string    `bson:"default_code,omitempty" json:"default_code"`
	Disabled    bool      `bson:"disabled" json:"disabled"`
	RegDate     time.Time `bson:"reg_date,omitempty" json:"reg_date"`
	UpdateDate  time.Time `bson:"update_date,omitempty" json:"update_date"`
}
