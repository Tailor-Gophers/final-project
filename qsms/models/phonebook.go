package models

type PhoneBook struct {
	Model
	UserID uint   `gorm:"not null"`
	Name   string `gorm:"size:255;not null"`
	Number Number `gorm:"foreignKey:PhoneBookID;references:ID"`
}
