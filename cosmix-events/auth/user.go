package auth

import "time"

type AuthUserRegistered struct {
	EventVersion string    `json:"event_version"`
	AuthUserID   uint      `json:"auth_user_id"`
	Email        string    `json:"email"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	CreatedAt    time.Time `json:"created_at"`
}