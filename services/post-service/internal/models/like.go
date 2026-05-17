package models

import "time"

type Like struct {
	ID        uint      `gorm:"primaryKey"`
	PostID    uint      `gorm:"not null;index"`
	UserID    uint      `gorm:"not null;index"`
	CreatedAt time.Time `gorm:"not null"`

	Post Post     `gorm:"foreignKey:PostID"`
	User PostUser `gorm:"foreignKey:UserID"`
}

func (Like) TableName() string {
	return "likes"
}
