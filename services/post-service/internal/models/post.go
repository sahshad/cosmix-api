package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	BaseModel
	UserID        uuid.UUID  `gorm:"type:uuid;not null;index"`
	LikesCount    int        `gorm:"default:0"`
	CommentsCount int        `gorm:"default:0"`
	Content       string     `gorm:"type:text;not null"`
	CreatedAt     time.Time  `gorm:"not null"`
	UpdatedAt     *time.Time `gorm:"default:null"`

	User     PostUser    `gorm:"foreignKey:UserID;references:UserID"`
	Media    []PostMedia `gorm:"foreignKey:PostID"`
	Likes    []Like      `gorm:"foreignKey:PostID"`
	Comments []Comment   `gorm:"foreignKey:PostID"`
}
