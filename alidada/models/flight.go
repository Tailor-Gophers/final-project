package models

import (
	"gorm.io/gorm"
	"time"
)

type Flight struct {
	gorm.Model
	Origin      string    `json:"flight_origin"`
	Destination string    `json:"flight_destination"`
	StartTime   time.Time `json:"flight_starttime"`
	EndTime     time.Time `json:"flight_endtime"`
	Airline     string    `json:"flight_airline"`
	Aircraft    string    `json:"flight_aircraft"`
	Capacity    uint      `json:"flight_capacity"`
}
