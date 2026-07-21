package dto

import (
	"time"

	"github.com/google/uuid"
)

type UpdateProfileDTO struct {
	DisplayName   *string
	Username      *string
	DateOfBirth   *string
	AvatarURL     *string
	CoverImageURL *string
	Website       *string
	Location      *string
	Bio           *string
}

type UserProfileResponse struct {
	User UserResponse
}

type UserResponse struct {
	UserID        uuid.UUID
	DisplayName   string
	Username      string
	Email         string
	IsPrivate     bool
	IsVerified    bool
	IsActive      bool
	DateOfBirth   *time.Time
	AvatarURL     *string
	CoverImageURL *string
	Website       *string
	Location      *string
	Bio           *string
	LastSeenAt    *time.Time
	CreatedAt     time.Time
	UpdatedAt     *time.Time
}

type UserCreatedFromDTO struct {
	AuthUserID uuid.UUID
	Email      string
	FirstName  string
	LastName   string
	CreatedAt  time.Time
}
