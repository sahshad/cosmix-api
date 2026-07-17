package models

import (
	"time"

	"github.com/google/uuid"
)

type Follow struct {
	BaseModel
	FollowerID  uuid.UUID `gorm:"not null;index:idx_follower_following,unique" json:"follower_id"`
	FollowingID uuid.UUID `gorm:"not null;index:idx_follower_following,unique" json:"following_id"`
	CreatedAt   time.Time `json:"created_at"`
}

func (Follow) TableName() string {
	return "follows"
}
