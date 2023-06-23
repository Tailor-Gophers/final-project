package models

import (
	"time"
)

type Token struct {
	Model
	UserId    uint
	Token     string `gorm:"unique;size:255"`
	ExpiresAt time.Time
}
