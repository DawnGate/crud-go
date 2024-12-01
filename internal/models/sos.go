package models

import (
	"time"

	"gorm.io/gorm"
)

type Sos struct {
	Title       string `json:"title" gorm:"not null"`
	Description string `json:"description"`

	// user relationship
	UserID uint `json:"userId"`
	User   User

	// orm default
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
