package models

type Aircraft struct {
	Id          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Type        string `gorm:"size:255;not null" json:"type"`
	Description string `gorm:"size:255" json:"description"`
}
