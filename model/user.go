package model

import "time"

type User struct {
	Id string `gorm:"primaryKey" json:"id" validate:"fixed"`

	Email          string `json:"email" validate:"required,fixed"`
	PasswordDigest string `json:"-"`

	FirstName string `json:"firstName"`
	LastName  string `json:"lastName" validate:"required"`

	CreatedAt time.Time `json:"created_at" validate:"zerotime"`
	UpdatedAt time.Time `json:"updated_at" validate:"zerotime"`
}
