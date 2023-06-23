package models

import "gorm.io/gorm"

type Reservation struct {
	gorm.Model
	PassengerID   uint
	FlightClassID uint
	OrderID       uint
	Price         uint
	IsCancelled   bool
}
