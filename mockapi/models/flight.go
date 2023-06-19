package models

import (
	"time"
)

type Flight struct {
	Id          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Origin      string    `gorm:"size:255" json:"origin"`
	Destination string    ` json:"destination"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Aircraft    Aircraft  `gorm:"foreignKey:Id" json:"aircraft"`
	Capacity    *uint     `gorm:"null" json:"capacity"`
}
