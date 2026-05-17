package models

import "time"

type Follow struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	FollowerID  uint      `gorm:"not null;index:idx_follower_following,unique" json:"follower_id"`
	FollowingID uint      `gorm:"not null;index:idx_follower_following,unique" json:"following_id"`
	CreatedAt   time.Time `json:"created_at"`
}

func (Follow) TableName() string {
	return "follows"
}
