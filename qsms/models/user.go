package models

type User struct {
	Model
	UserName     string      `gorm:"size:255;not null;unique" json:"user_name"`
	Email        string      `gorm:"size:255;not null;unique" json:"email"`
	Password     string      `gorm:"size:255;not null" json:"password"`
	Balance      int         `gorm:"default:0" json:"balance"`
	Disable      bool        `gorm:"default:false" json:"disable"`
	Admin        bool        `gorm:"default:false" json:"admin"`
	MainNumberID uint        `json:"main_number_id"`
	Numbers      []Number    `gorm:"foreignKey:ID" json:"numbers"`
	PhoneBooks   []PhoneBook `gorm:"foreignKey:ID" json:"phonebooks"`
	Contacts     []Contact   `json:"contacts"`
}
