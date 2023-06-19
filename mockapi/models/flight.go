package models

import (
	"time"
)

type Flight struct {
	ID          int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Origin      string `gorm:"size:255"`
	Destination string `gorm:"size:255"`
	StartTime   time.Time
	EndTime     time.Time
	Aircraft    string `gorm:"size:255"`
	Capacity    *uint  `gorm:"null"`
}
