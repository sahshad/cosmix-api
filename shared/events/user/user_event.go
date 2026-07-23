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
