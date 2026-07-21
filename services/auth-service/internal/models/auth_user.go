package models

import "time"

type AuthUser struct {
	BaseModel
	Email         string     `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash  string     `gorm:"type:varchar(255);not null"`
	IsActive      bool       `gorm:"not null;default:true"`
	DisplayName   string     `gorm:"type:varchar(100);not null"`
	EmailVerified bool       `gorm:"not null;default:false"`
	LastLoginAt   *time.Time `gorm:"default:null"`
	CreatedAt     time.Time  `gorm:"not null;default:now()"`
	UpdatedAt     *time.Time `gorm:"default:null"`

	AuthSessions []AuthSession `gorm:"foreignKey:AuthUserID"`
}

func (AuthUser) TableName() string {
	return "auth_users"
}
