package models

import "time"

type Comment struct {
	ID        uint       `gorm:"primaryKey"`
	PostID    uint       `gorm:"not null;index"`
	UserID    uint       `gorm:"not null;index"`
	Content   string     `gorm:"type:text;not null"`
	CreatedAt time.Time  `gorm:"not null"`
	UpdatedAt *time.Time `gorm:"default:null"`

	User PostUser `gorm:"foreignKey:UserID"`
	Post Post     `gorm:"foreignKey:PostID"`
}

func (Comment) TableName() string {
	return "comments"
}
