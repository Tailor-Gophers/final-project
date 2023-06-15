package models

import (
	"gorm.io/gorm"
	"time"
)

type AccessToken struct {
	gorm.Model
	Id        uint   `gorm:"primaryKey;autoIncrement"`
	UserId    uint   `gorm:"index"`
	Token     string `gorm:"uniqueIndex"`
	ExpiresAt time.Time
}
