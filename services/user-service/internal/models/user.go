package models

import "time"

type User struct {
	UserID      uint       `gorm:"primaryKey;not null"`
	Email       string     `gorm:"not null"`
	DisplayName string     `gorm:"not null"`
	Username    string     `gorm:"not null"`
	Bio         *string    `gorm:"default:null"`
	AvatarURL   *string    `gorm:"default:null"`
	IsPrivate   bool       `gorm:"default:false"`
	IsActive    bool       `gorm:"default:true"`
	DateOfBirth *time.Time `gorm:"default:null"`
	CreatedAt   time.Time  `gorm:"not null"`
	UpdatedAt   *time.Time `gorm:"default:null"`
}

func (User) TableName() string {
	return "users"
}
