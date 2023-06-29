package models

type Number struct {
	Model
	UserID      uint   `gorm:"not null"`
	UserName    string `json:"user_name"`
	PhoneNumber string `json:"phone_number"`
	Price       int    `json:"price"`
	Active      bool   `json:"active"`
	PhoneBookID uint   `gorm:"not null"`
}
