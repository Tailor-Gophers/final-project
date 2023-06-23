package models

type FlightClass struct {
	Model
	Title    string  `gorm:"size:255" json:"flight_class_title,omitempty"`
	Price    uint    `json:"flight_price,omitempty"`
	Capacity uint    `json:"flight_capacity,omitempty"`
	Reserve  *uint   `json:"flight_reserve,omitempty"`
	FlightId uint    `json:"-"`
	Flight   *Flight `json:"flight,omitempty"`
}
