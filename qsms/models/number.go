package models

type Number struct {
	Model
	UserID      uint   `gorm:"not null"`
	PhoneNumber string `json:"phone_number"`
	Price       int    `json:"price"`
	Active      bool   `json:"active"`
}
