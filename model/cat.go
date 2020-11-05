package model

import "time"

type Cat struct {
	Id     string `gorm:"primaryKey" json:"id"`
	UserId string `json:"UserId" validate:"fixed"`

	Name   string `json:"name"`
	Gender string `json:"gender" validate:"required,enum=MALE/FEMALE"`

	CreateTime time.Time `json:"createTime" validate:"zerotime"`
	UpdateTime time.Time `json:"updateTime" validate:"zerotime"`
}

func (c Cat) TableName() string {
	return "cats"
}
