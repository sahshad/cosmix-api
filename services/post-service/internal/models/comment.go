package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	BaseModel
	PostID    uuid.UUID  `gorm:"type:uuid;not null;index"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index"`
	Content   string     `gorm:"type:text;not null"`
	CreatedAt time.Time  `gorm:"not null"`
	UpdatedAt *time.Time `gorm:"default:null"`

	User PostUser `gorm:"foreignKey:UserID"`
	Post Post     `gorm:"foreignKey:PostID"`
}

func (Comment) TableName() string {
	return "comments"
}
