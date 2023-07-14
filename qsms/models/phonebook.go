package models

type PhoneBook struct {
	Model
	UserID  uint   `gorm:"not null" json:"user_id"`
	Name    string `gorm:"size:255;not null"`
	Numbers []RawNumber
}

type RawNumber struct {
	ID          uint `gorm:"primaryKey"`
	PhoneBookID uint
	Number      string
}
