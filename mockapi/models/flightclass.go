package models

import "gorm.io/gorm"

type FlightClass struct {
	gorm.Model
	Title    string `gorm:"size:255" json:"flight_class_title"`
	Price    uint   `json:"flight_price"`
	Capacity uint   `json:"flight_capacity"`
	Reserve  *uint  `json:"flight_reserve"`
	FlightId uint
	Flight   Flight
}
