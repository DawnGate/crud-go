package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique; not null"`
	Password string `json:"-" gorm:"not null"`

	// one-many relation
	Soss []Sos
}
