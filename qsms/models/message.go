package models

type Message struct {
	ID             uint   `gorm:"primaryKey"`
	SenderID       uint   `gorm:"not null" json:"sender_id"`
	ReceiverNumber string `gorm:"size:15" json:"receiver_number"`
	Message        string
}
