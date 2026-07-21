package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID        uuid.UUID  `gorm:"primaryKey;not null"`
	Email         string     `gorm:"not null"`
	Username      string     `gorm:"not null"`
	DisplayName   string     `gorm:"not null"`

	Bio           *string    `gorm:"default:null"`
	AvatarURL     *string    `gorm:"default:null"`
	CoverImageURL *string    `gorm:"default:null"`
	Website       *string    `gorm:"default:null"`
	Location      *string    `gorm:"default:null"`
	DateOfBirth   *time.Time `gorm:"default:null"`

	IsPrivate     bool       `gorm:"default:false"`
	IsVerified    bool       `gorm:"default:false"`
	IsActive      bool       `gorm:"default:true"`

	LastSeenAt    *time.Time `gorm:"default:null"`
	CreatedAt     time.Time  `gorm:"not null"`
	UpdatedAt     *time.Time `gorm:"default:null"`
}

func (User) TableName() string {
	return "users"
}
