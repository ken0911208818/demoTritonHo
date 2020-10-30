package model

import "time"

type Cat struct {
	Id string `gorm:"primaryKey" json:"id"`

	Name   string `json:"name"`
	Gender string `json:"gender"`

	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
}

func (c Cat) TableName() string {
	return "cats"
}
