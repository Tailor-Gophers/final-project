package models

type Template struct {
	Model
	UserID     uint   ` json:"user_id"`
	Expression string `gorm:"size:255" json:"expression"`
}

//Variables:
//	{{user_name}}, {{date}}, {{time}}, {{date_time}}
