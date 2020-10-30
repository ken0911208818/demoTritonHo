package handler

import (
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func Init(database *gorm.DB) {
	db = database

}
