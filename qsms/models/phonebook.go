package models

type PhoneBook struct {
	Model
	Name     string    `gorm:"size:255;not null"`
	Contacts []Contact `gorm:"foreignKey:PhoneBookID;references:ID"`
}
