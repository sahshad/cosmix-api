package user

import "time"

type UserUpdated struct {
	EventVersion string    `json:"event_version"`
	UserID       uint      `json:"user_id"`
	Username     string    `json:"username"`
	DisplayName  string    `json:"display_name"`
	AvatarURL    *string   `json:"avatar_url"`
	Bio          *string   `json:"bio"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserDeleted struct {
	EventVersion string `json:"event_version"`
	UserID       uint   `json:"user_id"`
}
