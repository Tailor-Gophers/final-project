package models

import "gorm.io/gorm"

type Reservation struct {
	gorm.Model
	PassengerID   uint         `json:"passenger_id,omitempty"`
	FlightClassID uint         `json:"flight_class_id,omitempty"`
	OrderID       uint         `json:"order_id,omitempty"`
	Price         uint         `json:"price,omitempty"`
	IsCancelled   bool         `json:"is_cancelled,omitempty"`
	FlightClass   *FlightClass `json:"flight_class,omitempty"`
	Passenger     *Passenger   `json:"passenger,omitempty"`
}
