package model

import "time"

type User struct {
	Id string `gorm:"primaryKey" json:"id" validate:"fixed"`

	Email          string `json:"email" validate:"required,fixed"`
	PasswordDigest string `json:"-"`

	FirstName string `json:"firstName"`
	LastName  string `json:"lastName" validate:"required"`

	CreateTime time.Time `json:"createTime" validate:"zerotime"`
	UpdateTime time.Time `json:"updateTime" validate:"zerotime"`
}
