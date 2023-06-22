package models

import "gorm.io/gorm"

type Reservation struct {
	gorm.Model
	PassengerID   uint
	FlightClassID uint
	Price         uint
	IsCancelled   bool
}
