package models

type Order struct {
	Model
	Reservations []Reservation
	Price        uint
	Confirmed    bool
}
