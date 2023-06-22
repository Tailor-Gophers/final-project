package models

import "gorm.io/gorm"

type Passenger struct {
	//Id             uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	gorm.Model
	UserID         uint
	FirstName      string `gorm:"size:255;not null" json:"first_name"`
	LastName       string `gorm:"size:255;not null" json:"last_name"`
	Gender         string `gorm:"size:255" json:"gender"`
	DateOfBirth    string `gorm:"size:255" json:"date_of_birth"`
	Nationality    string `gorm:"size:255" json:"nationality"`
	PassportNumber string `gorm:"size:255;not null" json:"passport_number"`
}
