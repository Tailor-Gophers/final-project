package models

type Reservation struct {
	Model
	PassengerID   uint
	FlightClassID uint
	OrderID       uint
	Price         uint
	IsCancelled   bool
	FlightData    FlightData `gorm:"foreignKey:ID"`
}
