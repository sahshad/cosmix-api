package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID      uuid.UUID `gorm:"type:uuid;primaryKey;not null"`
	Email       string    `gorm:"type:varchar(255);not null"`
	Username    string    `gorm:"type:varchar(255);not null"`
	DisplayName string    `gorm:"type:varchar(255);not null"`

	Bio           *string    `gorm:"type:text;default:null"`
	AvatarURL     *string    `gorm:"type:text;default:null"`
	CoverImageURL *string    `gorm:"type:text;default:null"`
	Website       *string    `gorm:"type:text;default:null"`
	Location      *string    `gorm:"type:text;default:null"`
	DateOfBirth   *time.Time `gorm:"type:date;default:null"`

	IsPrivate  bool `gorm:"default:false"`
	IsVerified bool `gorm:"default:false"`
	IsActive   bool `gorm:"default:true"`

	LastSeenAt *time.Time `gorm:"default:null"`
	CreatedAt  time.Time  `gorm:"not null"`
	UpdatedAt  *time.Time `gorm:"default:null"`
}

func (User) TableName() string {
	return "users"
}
