package models

import "time"

type AuthUser struct {
	ID            uint       `gorm:"primaryKey;not null"`
	Email         string     `gorm:"uniqueIndex;not null"`
	PasswordHash  string     `gorm:"not null"`
	IsActive      bool       `gorm:"default:true"`
	EmailVerified bool       `gorm:"default:false"`
	LastLoginAt   time.Time  `gorm:"not null"`
	CreatedAt     time.Time  `gorm:"not null"`
	UpdatedAt     *time.Time `gorm:"default:null"`
}

func (AuthUser) TableName() string {
	return "auth_users"
}