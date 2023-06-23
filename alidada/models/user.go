package models

type User struct {
	Model
	Email       string      `gorm:"not null;size:255" db:"email" json:"email"`
	UserName    string      `gorm:"not null;unique;size:255" db:"user_name" json:"user_name"`
	Password    string      `gorm:"not null;size:255" db:"password" json:"password"`
	FirstName   string      `gorm:"size:255" db:"first_name" json:"first_name"`
	LastName    string      `gorm:"size:255" db:"last_name" json:"last_name"`
	PhoneNumber string      `gorm:"size:11" json:"phone_number"`
	IsAdmin     bool        `gorm:"default:false" json:"is_admin"`
	Passengers  []Passenger `json:"passengers,omitempty"`
	Tokens      []Token     `json:"tokens,omitempty"`
}
