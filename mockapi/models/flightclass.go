package models

type FlightClass struct {
	Id       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Title    string `gorm:"size:255" json:"flight_class_title"`
	Price    uint   `json:"flight_price"`
	Capacity uint   `gorm:"null" json:"flight_capacity"`
	Reserve  *uint  `gorm:"null" json:"flight_reserve"`
}
