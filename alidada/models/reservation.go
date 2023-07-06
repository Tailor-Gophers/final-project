package models

type Reservation struct {
	Model
	PassengerID   uint
	FlightClassID uint
	OrderID       uint
	Price         uint
	IsCancelled   bool
	Passenger     Passenger
	Confirmed     bool        `gorm:"default:false"`
	FlightClass   FlightClass `gorm:"-:migration"`
}
