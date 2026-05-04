package dto

import "time"

type UpdateProfileDTO struct {
	Email       *string `json:"email"`
	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
	Username    *string `json:"username" binding:"omitempty,alphanum"`
	DateOfBirth *string `json:"date_of_birth"`
	AvatarURL   *string `json:"avatar_url"`
	Bio         *string `json:"bio" binding:"omitempty,max=500"`
}

type UserProfileResponse struct {
	User UserResponse `json:"user"`
}

type UserResponse struct {
	ID          uint       `json:"id"`
	AuthUserID  uint       `json:"auth_user_id"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Username    *string    `json:"username"`
	Email       *string    `json:"email"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	AvatarURL   *string    `json:"avatar_url,omitempty"`
	Bio         *string    `json:"bio,omitempty"`
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
