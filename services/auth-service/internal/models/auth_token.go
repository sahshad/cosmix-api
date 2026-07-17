package models

import (
	"time"

	"github.com/google/uuid"
)

type AuthTokenType string

const (
	AuthTokenTypeEmailVerification AuthTokenType = "email_verification"
	AuthTokenTypePasswordReset     AuthTokenType = "password_reset"
)

type AuthToken struct {
	BaseModel
	AuthUserID       uuid.UUID  `gorm:"not null"`
	Token            string     `gorm:"not null;unique"`
	Type             AuthTokenType `gorm:"not null"`
	ExpiresAt        time.Time  `gorm:"not null"`
	CreatedAt        time.Time  `gorm:"not null"`

	AuthUser AuthUser `gorm:"foreignKey:AuthUserID"`
}

func (AuthToken) TableName() string {
	return "auth_tokens"
}
