package models

import "time"

type PostMedia struct {
	ID        uint      `gorm:"primaryKey"`
	PostID    uint      `gorm:"not null;index"`
	PublicID  string    `gorm:"type:varchar(255);not null"`
	URL       string    `gorm:"type:text;not null"`
	Type      string    `gorm:"type:varchar(50);not null"`
	Duration  *int      `gorm:"default:null"`
	CreatedAt time.Time `gorm:"not null"`

	Post Post `gorm:"foreignKey:PostID"`
}

func (PostMedia) TableName() string {
	return "post_media"
}
