package models

import "gorm.io/gorm"

type FlightClass struct {
	gorm.Model
	Title    string `gorm:"size:255" json:"flight_class_title,omitempty"`
	Price    uint   `json:"flight_price,omitempty"`
	Capacity uint   `json:"flight_capacity,omitempty"`
	Reserve  *uint  `json:"flight_reserve,omitempty"`
	FlightId uint   `json:"flight_id,omitempty"`
	Flight   Flight `json:"flight,omitempty"`
}
