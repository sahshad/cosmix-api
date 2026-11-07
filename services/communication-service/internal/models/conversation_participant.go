package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	ParticipantRoleOwner  = "owner"
	ParticipantRoleAdmin  = "admin"
	ParticipantRoleMember = "member"
)

type ConversationParticipant struct {
	ConversationID uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID         uuid.UUID  `gorm:"type:uuid;primaryKey"`
	Role           string     `gorm:"type:varchar(20);not null;default:'member'"`
	JoinedAt       time.Time  `gorm:"not null"`
	LastReadAt     *time.Time `gorm:"default:null"`

	User CommunicationUser `gorm:"foreignKey:UserID;references:UserID"`
}

func (ConversationParticipant) TableName() string {
	return "conversation_participants"
}
