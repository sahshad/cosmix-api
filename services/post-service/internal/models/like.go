package models

import (
	"time"

	"github.com/google/uuid"
)

type Like struct {
	BaseModel
	PostID    uuid.UUID `gorm:"type:uuid;not null;index"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	CreatedAt time.Time `gorm:"not null"`

	Post Post     `gorm:"foreignKey:PostID"`
	User PostUser `gorm:"foreignKey:UserID"`
}

func (Like) TableName() string {
	return "likes"
}
