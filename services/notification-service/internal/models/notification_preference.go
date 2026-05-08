package models

import "time"

type NotificationPreference struct {
	ID                     uint       `gorm:"primaryKey" json:"id"`
	UserID                 uint       `gorm:"uniqueIndex;not null" json:"user_id"`
	EmailEnabled           bool       `gorm:"not null;default:true" json:"email_enabled"`
	PushEnabled            bool       `gorm:"not null;default:true" json:"push_enabled"`
	InternalEnabled        bool       `gorm:"not null;default:true" json:"internal_enabled"`
	LikeNotifications      bool       `gorm:"not null;default:true" json:"like_notifications"`
	CommentNotifications   bool       `gorm:"not null;default:true" json:"comment_notifications"`
	FollowNotifications    bool       `gorm:"not null;default:true" json:"follow_notifications"`
	MentionNotifications   bool       `gorm:"not null;default:true" json:"mention_notifications"`
	MessageNotifications   bool       `gorm:"not null;default:true" json:"message_notifications"`
	MarketingEmailsEnabled bool       `gorm:"not null;default:false" json:"marketing_emails_enabled"`
	CreatedAt              time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt              *time.Time `json:"updated_at"`
}

func (NotificationPreference) TableName() string {
	return "notification_preferences"
}
