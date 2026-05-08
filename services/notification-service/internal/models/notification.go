package models

import "time"

type Notification struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	UserID           uint       `gorm:"not null;index:idx_notifications_user_read" json:"user_id"`
	ActorID          *uint      `json:"actor_id"`
	ActorUsername    *string    `gorm:"size:255" json:"actor_username"`
	ActorDisplayName *string    `gorm:"size:255" json:"actor_display_name"`
	ActorAvatarURL   *string    `gorm:"type:text" json:"actor_avatar_url"`
	Type             string     `gorm:"size:100;not null" json:"type"`
	EntityID         *uint      `json:"entity_id"`
	EntityType       *string    `gorm:"size:100" json:"entity_type"`
	Title            string     `gorm:"type:text;not null" json:"title"`
	Body             string     `gorm:"type:text;not null" json:"body"`
	ImageURL         *string    `gorm:"type:text" json:"image_url"`
	ActionURL        *string    `gorm:"type:text" json:"action_url"`
	IsRead           bool       `gorm:"not null;default:false;index:idx_notifications_user_read" json:"is_read"`
	ReadAt           *time.Time `json:"read_at"`
	CreatedAt        time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
}

func (Notification) TableName() string {
	return "notifications"
}
