package models

type PhoneBook struct {
	Model
	UserID  uint      `gorm:"not null" json:"user_id"`
	Name    string    `gorm:"size:255;not null"`
	Numbers []*Number `gorm:"many2many:phone_book_number;"`
}
