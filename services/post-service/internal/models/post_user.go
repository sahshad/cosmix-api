package models

import (
	"time"

	"github.com/google/uuid"
)

type PostUser struct {
	UserID      uuid.UUID  `gorm:"type:uuid;primaryKey"`
	Username    string     `gorm:"type:varchar(100);not null;index"`
	DisplayName string     `gorm:"type:varchar(150);not null"`
	AvatarURL   *string    `gorm:"type:text;default:null"`
	IsVerified  bool       `gorm:"not null"`
	CreatedAt   time.Time  `gorm:"not null"`
	UpdatedAt   *time.Time `gorm:"default:null"`
}

func (PostUser) TableName() string {
	return "post_users"
}
