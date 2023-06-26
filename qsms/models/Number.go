package models

type Number struct {
	Model
	UserID      uint
	PhoneNumber string
	Price       int
	Active      bool
}

type Contact struct {
	Model
	UserID      uint   `gorm:"unique;index"`
	UserName    string `gorm:"size:255"`
	PhoneNumber string `gorm:"size:255"`
}
