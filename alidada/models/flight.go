package models

import (
	"time"

	"gorm.io/gorm"
)

type Flight struct {
	Model
	Origin      string    `gorm:"size:255" json:"flight_origin,omitempty"`
	Destination string    `gorm:"size:255" json:"flight_destination,omitempty"`
	StartTime   time.Time `json:"flight_starttime,omitempty"`
	EndTime     time.Time `json:"flight_endtime,omitempty"`
	Airline     string    `gorm:"size:255" json:"flight_airline,omitempty"`
	Aircraft    string    `gorm:"size:255" json:"flight_aircraft,omitempty"`
}

type Model struct {
	ID        uint           `gorm:"primaryKey"`
	CreatedAt time.Time      `json:"omitempty"`
	UpdatedAt time.Time      `json:"omitempty"`
	DeletedAt gorm.DeletedAt `json:"omitempty" gorm:"index"`
}
