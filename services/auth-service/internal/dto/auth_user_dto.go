package dto

type AuthUserDTO struct {
	Email         string
	IsActive      bool
	EmailVerified bool
	LastLoginAt   string
	CreatedAt     string
	UpdatedAt     string
}
