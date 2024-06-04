package models

import (
	"time"

	"github.com/sev-2/raiden"
)

type List struct {
	raiden.ModelBase
	Id         int64     `json:"id,omitempty" column:"name:id;type:bigint;primaryKey;autoIncrement;nullable:false"`
	CreatedAt  time.Time `json:"created_at,omitempty" column:"name:created_at;type:timestampz;nullable:false;default:now()"`
	Text       string    `json:"text,omitempty" column:"name:text;type:text;nullable:false;default:null"`
	IsComplete bool      `json:"is_complete,omitempty" column:"name:is_complete;type:bool;nullable:false;default:false"`

	// Table information
	Metadata string `json:"-" schema:"public" rlsEnable:"false" rlsForced:"false"`

	// Access control
	Acl string `json:"-" read:"" write:""`
}
