package Models

import "time"

type Message struct {
	ID int
	SenderID int
	Sender User `gorm:"foreignKey:SenderID"`
	RecipientID int
	Recipient User `gorm:"foreignKey:RecipientID"`
	CreatedAt time.Time
	Body string
}

