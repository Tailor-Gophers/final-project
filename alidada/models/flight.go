package models

import "time"

type Flight struct {
	Id          int64     `json:"flight_id"`
	Origin      string    `json:"flight_origin"`
	Destination string    `json:"flight_destination"`
	StartTime   time.Time `json:"flight_starttime"`
	EndTime     time.Time `json:"flight_endtime"`
	Airline     string    `json:"flight_airline"`
	Aircraft    string    `json:"flight_aircraft"`
	Reserve     uint      `json:"flight_reserve"`
}
