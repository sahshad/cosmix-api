package models

import "time"

type UserProfile struct {
	ID          uint       `gorm:"primaryKey"`
	AuthUserID  uint       `gorm:"uniqueIndex;not null"`
	Email       string     `gorm:"uniqueIndex;not null"`
	FirstName   string     `gorm:"not null"`
	LastName    string     `gorm:"not null"`
	Username    *string    `gorm:"uniqueIndex"`
	DisplayName *string    `json:"display_name"`
	Bio         *string    `json:"bio"`
	AvatarURL   *string    `json:"avatar_url"`
	IsPrivate   bool       `gorm:"default:false"`
	IsActive    bool       `gorm:"default:true"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	CreatedAt   time.Time  `json:"created_at" gorm:"not null"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"default:null"`
}
