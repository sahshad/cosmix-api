package dto

import (
	"time"

	"github.com/google/uuid"
)

type UpdateProfileDTO struct {
	DisplayName   *string `json:"display_name"`
	Username      *string `json:"username"`
	DateOfBirth   *string `json:"date_of_birth"`
	AvatarURL     *string `json:"avatar_url"`
	CoverImageURL *string `json:"cover_image_url"`
	Website       *string `json:"website" binding:"max=255"`
	Location      *string `json:"location" binding:"max=255"`
	Bio           *string `json:"bio" binding:"max=500"`
}

type UserProfileResponse struct {
	User UserResponse `json:"user"`
}

type UserResponse struct {
	UserID        uuid.UUID  `json:"id"`
	DisplayName   string     `json:"display_name"`
	Username      string     `json:"username"`
	Email         string     `json:"email"`
	IsPrivate     bool       `json:"is_private"`
	IsVerified    bool       `json:"is_verified"`
	IsActive      bool       `json:"is_active"`
	DateOfBirth   *time.Time `json:"date_of_birth"`
	AvatarURL     *string    `json:"avatar_url"`
	CoverImageURL *string    `json:"cover_image_url"`
	Website       *string    `json:"website"`
	Location      *string    `json:"location"`
	Bio           *string    `json:"bio"`
	LastSeenAt    *time.Time `json:"last_seen_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
}

type UserCreatedFromDTO struct {
	AuthUserID uuid.UUID
	Email      string
	FirstName  string
	LastName   string
	CreatedAt  time.Time
}
