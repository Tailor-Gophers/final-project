package models

type Transaction struct {
	Model
	UserID    uint
	Authority string `gorm:"size:255"`
	Amount    int
	RefID     int
	Confirmed bool
}
