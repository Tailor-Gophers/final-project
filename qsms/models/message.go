package models

type Message struct {
	Model
	SenderID       uint   `gorm:"not null" json:"sender_id"`
	ReceiverNumber string `gorm:"size:12" json:"receiver_number"`
	Message        string
}
