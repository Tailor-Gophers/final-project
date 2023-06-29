package models

type PhoneBook struct {
	Model
	ID      uint     `gorm:"primaryKey;autoIncrement"`
	UserID  uint     `gorm:"not null;foreignKey:UserID" json:"user_id"`
	Name    string   `gorm:"size:255;not null"`
	Numbers []Number `gorm:"foreignKey:ID"`
}
