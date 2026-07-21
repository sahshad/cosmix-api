package dto

import (
	"time"

	"github.com/google/uuid"
)

type RegisterDTO struct {
	DisplayName string
	Email       string
	Password    string
}

type VerifyEmailDTO struct {
	Token    string
	Email    string
	Password string
}

type LoginDTO struct {
	Email    string
	Password string
}

type UserUpdatedFromDTO struct {
	AuthUserID uuid.UUID
	Email      string
	UpdatedAt  time.Time
}

type UpdateUserPasswordDTO struct {
	UserID      uuid.UUID
	NewPassword string
}

type LoginResponseDTO struct {
	AccessToken  string
	RefreshToken string
	AuthUser     *AuthUserDTO
}

type RefreshResponseDTO struct {
	AccessToken  string
	RefreshToken string
}

type ResetPasswordDTO struct {
	Token       string
	NewPassword string
}
