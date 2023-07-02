package models

import "time"

// dont nead
type FlightData struct {
	Model
	Origin      string    `gorm:"size:255" json:"origin,omitempty"`
	Destination string    `gorm:"size:255" json:"destination,omitempty"`
	Title       string    `gorm:"size:255" json:"class_title,omitempty"`
	Airline     string    `gorm:"size:255" json:"airline,omitempty"`
	Aircraft    string    `gorm:"size:255" json:"aircraft,omitempty"`
	StartTime   time.Time `json:"starttime,omitempty"`
	EndTime     time.Time `json:"endtime,omitempty"`
}
