package models

import "time"

type Order struct {
	Model
	Reservations []Reservation
	Price        uint
	OrderTime    time.Time
	RefID        int
	Confirmed    bool `gorm:"default:false"`
}
