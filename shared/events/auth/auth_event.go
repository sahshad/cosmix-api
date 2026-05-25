package auth

import "time"

const (
	EventVersionOne = "1.0.0"
)

type AuthUserRegistered struct {
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
