package models

type FlightClass struct {
	Id       int64  `json:"id"`
	Title    string `json:"flight_class_title"`
	Price    uint   `json:"flight_price"`
	Capacity uint   `json:"flight_capacity"`
	Reserve  *uint  `json:"flight_reserve"`
	FlightId uint
	Flight   Flight
}
