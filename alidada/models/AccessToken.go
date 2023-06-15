package models

import (
	"time"

	"gorm.io/gorm"
)

type AccessToken struct {
	gorm.Model
	UserId    uint
	Token     string `gorm:"unique;size:255"`
	ExpiresAt time.Time
}
