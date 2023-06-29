package models

type Contact struct {
	Model
	UserID      uint
	Name        string `gorm:"size:255" json:"name"`
	PhoneNumber string `gorm:"size:255" json:"phone"`
}
