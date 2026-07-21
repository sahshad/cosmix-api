package models

import (
	"time"

	"github.com/google/uuid"
)

type Like struct {
	BaseModel
	PostID    uuid.UUID `gorm:"type:uuid;not null;index:idx_likes_post_user,unique"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index:idx_likes_post_user,unique"`
	CreatedAt time.Time `gorm:"not null"`

	Post Post     `gorm:"foreignKey:PostID"`
	User PostUser `gorm:"foreignKey:UserID"`
}

func (Like) TableName() string {
	return "likes"
}
