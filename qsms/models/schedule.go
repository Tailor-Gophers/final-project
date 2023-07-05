package models

import (
	"time"
)

type Rent struct {
	Model
	UserID   uint
	NumberID uint
	Price    int
	LastPaid time.Time
}

type MessageSchedule struct {
	Model
	UserID   uint
	Receiver string
	Text     string
	Template string
	Interval string
}
