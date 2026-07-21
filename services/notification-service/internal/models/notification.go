package models

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	BaseModel
	UserID           uuid.UUID  `gorm:"type:uuid;not null;index:idx_notifications_user_read"`
	ActorID          *uuid.UUID `gorm:"type:uuid"`
	ActorUsername    *string    `gorm:"size:255"`
	ActorDisplayName *string    `gorm:"size:255"`
	ActorAvatarURL   *string    `gorm:"type:text"`
	Type             string     `gorm:"size:100;not null"`
	EntityID         *uuid.UUID `gorm:"type:uuid"`
	EntityType       *string    `gorm:"size:100"`
	Title            string     `gorm:"type:text;not null"`
	Body             string     `gorm:"type:text;not null"`
	ImageURL         *string    `gorm:"type:text"`
	ActionURL        *string    `gorm:"type:text"`
	IsRead           bool       `gorm:"not null;default:false;index:idx_notifications_user_read"`
	ReadAt           *time.Time `gorm:"default:null"`
	CreatedAt        time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func (Notification) TableName() string {
	return "notifications"
}
