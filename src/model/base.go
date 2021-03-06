package model

import (
	"encoding/json"
	"time"

	"github.com/moocss/apiserver/src/util"
)

type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt util.JSONTime
	UpdatedAt util.JSONTime
	DeletedAt *util.JSONTime `sql:"index"`
}

type BaseModel struct {
	ID        uint64     `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"-"`
	CreatedAt time.Time  `gorm:"column:createdAt" json:"-"`
	UpdatedAt time.Time  `gorm:"column:updatedAt" json:"-"`
	DeletedAt *time.Time `gorm:"column:deletedAt" sql:"index" json:"-"`
}

func jsonMarshal(v interface{}) (str string) {
	if res, err := json.Marshal(v); err == nil {
		str = string(res)
	}
	return
}
