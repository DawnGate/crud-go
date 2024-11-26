package models

import (
	"time"

	"gorm.io/gorm"
)

type Sos struct {
	ID          uint `gorm:"primaryKey"`
	Title       string
	Description string
	UserId      uint
	User        User
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}
