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
	Revoked          bool       `gorm:"default:false"`
	CreatedAt        time.Time  `gorm:"not null"`
	UpdatedAt        *time.Time `gorm:"default:null"`

	// Relation
	AuthUser AuthUser `gorm:"foreignKey:AuthUserID;references:ID"`
}

func (AuthSession) TableName() string {
	return "auth_sessions"
}
