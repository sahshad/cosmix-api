package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	BaseModel
	ConversationID uuid.UUID  `gorm:"type:uuid;not null;index"`
	SenderID       uuid.UUID  `gorm:"type:uuid;not null;index"`
	Content        string     `gorm:"type:text;not null"`
	CreatedAt      time.Time  `gorm:"not null"`
	UpdatedAt      *time.Time `gorm:"default:null"`

	Sender CommunicationUser `gorm:"foreignKey:SenderID;references:UserID"`
}

func (Message) TableName() string {
	return "messages"
}
