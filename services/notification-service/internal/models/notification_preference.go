package models

import (
	"time"

	"github.com/google/uuid"
)

type NotificationPreference struct {
	BaseModel
	UserID                 uuid.UUID  `gorm:"type:uuid;uniqueIndex;not null"`
	EmailEnabled           bool       `gorm:"not null;default:true"`
	PushEnabled            bool       `gorm:"not null;default:true"`
	InternalEnabled        bool       `gorm:"not null;default:true"`
	LikeNotifications      bool       `gorm:"not null;default:true"`
	CommentNotifications   bool       `gorm:"not null;default:true"`
	FollowNotifications    bool       `gorm:"not null;default:true"`
	MentionNotifications   bool       `gorm:"not null;default:true"`
	MessageNotifications   bool       `gorm:"not null;default:true"`
	MarketingEmailsEnabled bool       `gorm:"not null;default:false"`
	CreatedAt              time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt              *time.Time `gorm:"default:null"`
}

func (NotificationPreference) TableName() string {
	return "notification_preferences"
}
