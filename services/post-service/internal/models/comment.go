package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	BaseModel
	PostID          uuid.UUID  `gorm:"type:uuid;not null;index"`
	UserID          uuid.UUID  `gorm:"type:uuid;not null;index"`
	ParentCommentID *uuid.UUID `gorm:"type:uuid;index"`
	Content         string     `gorm:"type:text;not null"`
	LikesCount      int        `gorm:"default:0"`
	RepliesCount    int        `gorm:"default:0"`
	CreatedAt       time.Time  `gorm:"not null"`
	UpdatedAt       *time.Time `gorm:"default:null"`

	User PostUser `gorm:"foreignKey:UserID;references:UserID"`
	Post Post     `gorm:"foreignKey:PostID"`
}

func (Comment) TableName() string {
	return "comments"
}
