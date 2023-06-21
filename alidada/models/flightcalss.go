package models

type FlightClass struct {
	Id       int64  `json:"class_id"`
	Title    string `json:"class_title"`
	Capacity uint   `json:"flight_capacity"`
	Price    uint   `json:"flight_price"`
	FlightId int64
	Flight   Flight
}
