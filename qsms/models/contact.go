package models

type Contact struct {
	Model
	PhoneBookID uint   `gorm:"not null"`
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"size:255"`
	PhoneNumber string `gorm:"size:255" json:"phone"`
}
