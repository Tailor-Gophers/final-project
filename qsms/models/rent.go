package models

import "time"

type Rent struct {
	Model
	UserID   uint
	NumberID uint
	Price    int
	LastPaid time.Time
}
