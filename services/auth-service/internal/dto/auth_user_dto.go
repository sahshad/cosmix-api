package dto

import "time"

type AuthUserDTO struct {
	Email         string
	IsActive      bool
	EmailVerified bool
	LastLoginAt   *time.Time
	CreatedAt     time.Time
	UpdatedAt     *time.Time
}
