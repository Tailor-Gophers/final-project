package models

type Number struct {
	Model
	UserName    string `json:"user_name"`
	PhoneNumber string `json:"phone_number"`
	Price       int    `json:"price"`
	Active      bool   `json:"active"`
}
