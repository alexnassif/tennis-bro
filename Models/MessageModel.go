package Models

import (
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	SenderID int
	Sender User `gorm:"foreignKey:SenderID"`
	RecipientID int
	Recipient User `gorm:"foreignKey:RecipientID"`
	Body string
}

