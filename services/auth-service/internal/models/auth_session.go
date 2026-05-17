package models

import "time"

type AuthSession struct {
	ID               uint       `gorm:"primaryKey;not null"`
	AuthUserID       uint       `gorm:"not null"`
	RefreshTokenHash string     `gorm:"not null"`
	Device           string     `gorm:"not null"`
	IPAddress        string     `gorm:"not null"`
	UserAgent        string     `gorm:"not null"`
	ExpiresAt        time.Time  `gorm:"not null"`
	Revoked          bool       `gorm:"not null"`
	CreatedAt        time.Time  `gorm:"not null"`
	UpdatedAt        *time.Time `gorm:"null"`

	AuthUser AuthUser `gorm:"foreignKey:AuthUserID"`
}

func (AuthSession) TableName() string {
	return "auth_sessions"
}
