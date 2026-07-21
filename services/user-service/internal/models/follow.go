package models

import (
	"time"

	"github.com/google/uuid"
)

type Follow struct {
	BaseModel
	FollowerID  uuid.UUID `gorm:"type:uuid;not null;index:idx_follower_following,unique"`
	FollowingID uuid.UUID `gorm:"type:uuid;not null;index:idx_follower_following,unique"`
	CreatedAt   time.Time `gorm:"not null"`
}

func (Follow) TableName() string {
	return "follows"
}
