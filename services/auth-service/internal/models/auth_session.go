package models

import (
	"time"

	"github.com/google/uuid"
)

type AuthSession struct {
	BaseModel
	AuthUserID       uuid.UUID  `gorm:"type:uuid;not null"`
	RefreshTokenHash string     `gorm:"type:varchar(500);not null"`
	Device           string     `gorm:"type:varchar(255)"`
	IPAddress        string     `gorm:"type:varchar(100)"`
	UserAgent        string     `gorm:"type:varchar(500)"`
	ExpiresAt        time.Time  `gorm:"not null"`
	Revoked          bool       `gorm:"not null"`
	CreatedAt        time.Time  `gorm:"not null"`
	UpdatedAt        *time.Time `gorm:"null"`

	AuthUser AuthUser `gorm:"foreignKey:AuthUserID"`
}

func (AuthSession) TableName() string {
	return "auth_sessions"
}
