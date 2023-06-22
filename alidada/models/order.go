package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	Reservations []Reservation
	Price        uint
	Confirmed    bool
}
