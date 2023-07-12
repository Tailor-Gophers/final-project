package models

type Number struct {
	Model
	UserID      *uint  `json:"user_id"`
	PhoneNumber string `gorm:"size:255;not null" json:"phone_number"`
	Price       int    `json:"price"`
	Active      bool   `json:"active"`
	//PhoneBook   []*PhoneBook `gorm:"many2many:phone_book_number;"`
}
