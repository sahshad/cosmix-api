package models

import (
	"time"

	"github.com/google/uuid"
)

type CommentLike struct {
	BaseModel
	CommentID uuid.UUID `gorm:"type:uuid;not null;index:idx_comment_likes_comment_user,unique"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index:idx_comment_likes_comment_user,unique"`
	CreatedAt time.Time `gorm:"not null"`

	Comment Comment  `gorm:"foreignKey:CommentID"`
	User    PostUser `gorm:"foreignKey:UserID;references:UserID"`
}

func (CommentLike) TableName() string {
	return "comment_likes"
}
