package user

import (
	"time"

	"github.com/google/uuid"
)

const (
	EventVersionOne = "1.0.0"
)

type UserUpdated struct {
	EventVersion string
	UserID       uuid.UUID
	Username     string
	DisplayName  string
	AvatarURL    *string
	IsVerified   bool
	UpdatedAt    time.Time
}

type UserDeleted struct {
	EventVersion string
	UserID       uuid.UUID
}

type UserFollowed struct {
	EventVersion      string
	FollowerID        uuid.UUID
	FollowerUsername  string
	FollowerDisplayName string
	FollowingID       uuid.UUID
	CreatedAt         time.Time
}
