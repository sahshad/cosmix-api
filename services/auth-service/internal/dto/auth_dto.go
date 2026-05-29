package dto

import "time"

type RegisterDTO struct {
	DisplayName string `json:"display_name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
}

type VerifyEmailDTO struct {
	Token    string `json:"token" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type LoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserUpdatedFromDTO struct {
	AuthUserID uint
	Email      string
	UpdatedAt  time.Time
}

type UpdateUserPasswordDTO struct {
	UserID      uint   `json:"user_id" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type LoginResponseDTO struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	AuthUser     *AuthUserDTO `json:"auth_user"`
}

type RefreshResponseDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ResetPasswordDTO struct {
	Token           string `json:"token" binding:"required"`
	CurrentPassword string `json:"current_password" binding:"required,min=8,max=64"`
	NewPassword     string `json:"new_password" binding:"required,min=8,max=64"`
}
