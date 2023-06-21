package models

import (
	"time"
)

type Flight struct {
	Id          int64     `gorm:"primaryKey;autoIncrement" json:"flight_id"`
	Origin      string    `gorm:"size:255" json:"flight_origin"`
	Destination string    `json:"flight_destination"`
	StartTime   time.Time `json:"flight_starttime"`
	EndTime     time.Time `json:"flight_endtime"`
	Airline     string    `gorm:"size:255" json:"flight_airline"`
	Aircraft    string    `gorm:"size:255" json:"flight_aircraft"`
}
