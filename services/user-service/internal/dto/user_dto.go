package dto

import "time"

type UpdateProfileDTO struct {
	DisplayName *string `json:"display_name"`
	Username    *string `json:"username"`
	DateOfBirth *string `json:"date_of_birth"`
	AvatarURL   *string `json:"avatar_url"`
	Bio         *string `json:"bio" binding:"max=500"`
}

type UserProfileResponse struct {
	User UserResponse `json:"user"`
}

type UserResponse struct {
	UserID      uint       `json:"id"`
	DisplayName string     `json:"display_name"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	IsPrivate   bool       `json:"is_private"`
	IsActive    bool       `json:"is_active"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	AvatarURL   *string    `json:"avatar_url"`
	Bio         *string    `json:"bio"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type UserCreatedFromDTO struct {
	AuthUserID uint
	Email      string
	FirstName  string
	LastName   string
	CreatedAt  time.Time
}
