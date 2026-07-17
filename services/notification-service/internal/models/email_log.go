package models

import (
	"time"

	"github.com/google/uuid"
)

type EmailLog struct {
	BaseModel
	UserID       *uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	Recipient    string     `gorm:"size:255;not null" json:"recipient"`
	Type         string     `gorm:"size:100;not null" json:"type"`
	Subject      string     `gorm:"type:text;not null" json:"subject"`
	Template     string     `gorm:"size:255;not null" json:"template"`
	Status       string     `gorm:"size:50;not null;index" json:"status"`
	Provider     *string    `gorm:"size:100" json:"provider"`
	ErrorMessage *string    `gorm:"type:text" json:"error_message"`
	SentAt       *time.Time `json:"sent_at"`
	FailedAt     *time.Time `json:"failed_at"`
	CreatedAt    time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (EmailLog) TableName() string {
	return "email_logs"
}
