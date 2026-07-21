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
	AuthUserID uuid.UUID     `gorm:"type:uuid;not null"`
	Token      string        `gorm:"type:varchar(255);not null;unique"`
	Type       AuthTokenType `gorm:"type:varchar(50);not null;check:type IN ('email_verification','password_reset')"`
	ExpiresAt  time.Time     `gorm:"not null"`
	CreatedAt  time.Time     `gorm:"not null"`

	AuthUser AuthUser `gorm:"foreignKey:AuthUserID"`
}

func (AuthToken) TableName() string {
	return "auth_tokens"
}
