package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName string `json:"user_name" gorm:"unique; not null"`
	Password string `json:"password" gorm:"not null"`
}
