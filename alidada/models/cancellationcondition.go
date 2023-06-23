package models

type CancellationCondition struct {
	ID          uint
	TimeMinutes uint   `json:"-"`
	Title       string `gorm:"size:255" json:"title,omitempty"`
	Description string `gorm:"size:255" json:"description,omitempty"`
	Penalty     int    `json:"penalty,omitempty"`
}
