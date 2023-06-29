package models

type PhoneBook struct {
	Model
	UserID  uint     `gorm:"not null"`
	Name    string   `gorm:"size:255;not null"`
	Numbers []Number `gorm:"foreignKey:ID"`
}
