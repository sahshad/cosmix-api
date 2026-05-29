package models

import "time"

type EmailVerificationToken struct {
	ID          uint      `gorm:"primaryKey;not null"`
	AuthUserID  uint      `gorm:"not null"`
	DisplayName string    `gorm:"not null"`
	Token       string    `gorm:"not null;unique"`
	ExpiresAt   time.Time `gorm:"not null"`
	CreatedAt   time.Time `gorm:"not null"`

	AuthUser AuthUser `gorm:"foreignKey:AuthUserID"`
}

func (EmailVerificationToken) TableName() string {
	return "email_verification_tokens"
}
