package models

import (
	"time"

	"gorm.io/gorm"
)

type Token struct {
	gorm.Model
	UserId    uint
	Token     string `gorm:"unique;size:255"`
	ExpiresAt time.Time
}
