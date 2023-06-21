package models

type FlightClass struct {
	Id       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Title    string `gorm:"size:255"`
	Capacity uint   `gorm:"null" json:"flight_capacity"`
	Price    uint   `json:"flight_price"`
	FlightId int64
	Flight   Flight
}
