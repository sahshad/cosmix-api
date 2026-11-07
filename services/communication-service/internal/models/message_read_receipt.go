package models

import (
	"time"

	"github.com/google/uuid"
)

type MessageReadReceipt struct {
	MessageID      uuid.UUID `gorm:"type:uuid;primaryKey"`
	ConversationID uuid.UUID `gorm:"type:uuid;not null;index"`
	UserID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	ReadAt         time.Time `gorm:"not null"`

	User CommunicationUser `gorm:"foreignKey:UserID;references:UserID"`
}

func (MessageReadReceipt) TableName() string {
	return "message_read_receipts"
}
