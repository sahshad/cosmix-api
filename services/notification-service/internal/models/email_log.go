package models

import (
	"time"

	"github.com/google/uuid"
)

type EmailLog struct {
	BaseModel
	UserID       *uuid.UUID `gorm:"type:uuid;index"`
	Recipient    string     `gorm:"size:255;not null"`
	Type         string     `gorm:"size:100;not null"`
	Subject      string     `gorm:"type:text;not null"`
	Template     string     `gorm:"size:255;not null"`
	Status       string     `gorm:"size:50;not null;index"`
	Provider     *string    `gorm:"size:100"`
	ErrorMessage *string    `gorm:"type:text"`
	SentAt       *time.Time `gorm:"default:null"`
	FailedAt     *time.Time `gorm:"default:null"`
	CreatedAt    time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func (EmailLog) TableName() string {
	return "email_logs"
}
