package auth

import "time"

const (
	EventVersionOne = "1.0.0"
)

type AuthUserEmailVerificationCompleted struct {
	EventVersion string    `json:"event_version"`
	AuthUserID   uint      `json:"auth_user_id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	DisplayName  string    `json:"display_name"`
	CreatedAt    time.Time `json:"created_at"`
}

type AuthUserEmailUpdated struct {
	EventVersion string `json:"event_version"`
	Email        string `json:"email"`
}

type AuthUserEmailVerification struct {
	EventVersion string `json:"event_version"`
	Email        string `json:"email"`
	DisplayName  string `json:"dispaly_name"`
	Token        string `json:"token"`
}

type AuthUserForgotPasswordRequest struct {
	EventVersion string `json:"event_version"`
	Email        string `json:"email"`
	DisplayName  string `json:"dispaly_name"`
	Token        string `json:"token"`
}

type AuthUserPasswordChanged struct {
	EventVersion string    `json:"event_version"`
	AuthUserID   uint      `json:"auth_user_id"`
	Email        string    `json:"email"`
	DisplayName  string    `json:"display_name"`
	UpdatedAt    time.Time `json:"updated_at"`
}
