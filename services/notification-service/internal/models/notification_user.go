package models

import "time"

type NotificationUser struct {
	UserID      uint       `gorm:"primaryKey"`
	Username    string     `gorm:"type:varchar(100);not null"`
	DisplayName string     `gorm:"type:varchar(150);not null"`
	AvatarURL   *string    `gorm:"type:text;default:null"`
	IsVerified  bool       `gorm:"not null"`
	CreatedAt   time.Time  `gorm:"not null"`
	UpdatedAt   *time.Time `gorm:"default:null"`
}

func (NotificationUser) TableName() string {
	return "notification_users"
}
