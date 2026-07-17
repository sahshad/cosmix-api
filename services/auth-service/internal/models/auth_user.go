package models

import "time"

type AuthUser struct {
	BaseModel
	Email         string     `gorm:"uniqueIndex; not null"`
	PasswordHash  string     `gorm:"not null"`
	IsActive      bool       `gorm:"not null;default:true"`
	DisplayName   string     `gorm:"not null"`
	EmailVerified bool       `gorm:"not null;default:false"`
	LastLoginAt   *time.Time `gorm:"default:null"`
	CreatedAt     time.Time  `gorm:"not null;default:now()"`
	UpdatedAt     *time.Time `gorm:"default:null"`

	AuthSessions []AuthSession `gorm:"foreignKey:AuthUserID"`
}

func (AuthUser) TableName() string {
	return "auth_users"
}
