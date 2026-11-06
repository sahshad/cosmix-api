package models

import (
	"time"

	"github.com/google/uuid"
)

type ConversationParticipant struct {
	ConversationID uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID         uuid.UUID  `gorm:"type:uuid;primaryKey"`
	JoinedAt       time.Time  `gorm:"not null"`
	LastReadAt     *time.Time `gorm:"default:null"`

	User CommunicationUser `gorm:"foreignKey:UserID;references:UserID"`
}

func (ConversationParticipant) TableName() string {
	return "conversation_participants"
}
