package models

type Passenger struct {
	//Id             uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Model
	UserID         uint   `json:"-"`
	FirstName      string `gorm:"size:255;not null" json:"first_name,omitempty"`
	LastName       string `gorm:"size:255;not null" json:"last_name,omitempty"`
	Gender         bool   `json:"gender,omitempty"`
	DateOfBirth    string `gorm:"size:255" json:"date_of_birth,omitempty"`
	Nationality    string `gorm:"size:10" json:"nationality,omitempty"`
	PassportNumber string `gorm:"size:10" json:"passport_number,omitempty"`
}
