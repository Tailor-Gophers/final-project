package models

type Contact struct {
	Model
	UserID      uint   `gorm:"not null"`
	Name        string `gorm:"size:255" json:"name"`
	PhoneNumber string `gorm:"size:255" json:"phone"`
}
