package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName    string `gorm:"not null;unique" db:"user_name" json:"user_name"`
	Password    string `gorm:"not null" db:"password" json:"password"`
	FirstName   string `db:"first_name" json:"first_name"`
	LastName    string `db:"last_name" json:"last_name"`
	PhoneNumber string `gorm:"size:10" json:"phone_number"`
	Admin       bool   `db:"admin" json:"admin"`
}
