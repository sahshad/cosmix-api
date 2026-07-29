package dto

import (
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
	UserID         uuid.UUID
	DisplayName    string
	Username       string
	Email          string
	IsPrivate      bool
	IsVerified     bool
	IsActive       bool
	DateOfBirth    string
	AvatarURL      *string
	CoverImageURL  *string
	Website        *string
	Location       *string
	Bio            *string
	LastSeenAt     *string
	FollowersCount int
	FollowingCount int
	IsFollowing    bool
	CreatedAt      string
	UpdatedAt      string
}

type UserCreatedFromDTO struct {
	AuthUserID uuid.UUID
	Email      string
	FirstName  string
	LastName   string
	CreatedAt  string
}
