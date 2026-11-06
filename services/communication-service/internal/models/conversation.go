package models

import (
	"time"

	"github.com/google/uuid"
)

const ConversationTypeDirect = "direct"

type Conversation struct {
	BaseModel
	Type          string     `gorm:"type:varchar(20);not null;default:'direct'"`
	DirectKey     *string    `gorm:"type:varchar(80);default:null"`
	CreatedBy     uuid.UUID  `gorm:"type:uuid;not null"`
	LastMessageID *uuid.UUID `gorm:"type:uuid;default:null"`
	LastMessageAt *time.Time `gorm:"default:null"`
	CreatedAt     time.Time  `gorm:"not null"`
	UpdatedAt     *time.Time `gorm:"default:null"`

	Participants []ConversationParticipant `gorm:"foreignKey:ConversationID"`
	LastMessage  *Message                  `gorm:"foreignKey:LastMessageID"`
}

func (Conversation) TableName() string {
	return "conversations"
}
