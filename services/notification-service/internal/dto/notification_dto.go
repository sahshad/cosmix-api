package dto

import (
	"time"

	"github.com/google/uuid"
)

type NotificationList struct {
	ID               uuid.UUID  `json:"id"`
	UserID           uuid.UUID  `json:"user_id"`
	ActorID          *uuid.UUID `json:"actor_id"`
	ActorUsername    *string    `json:"actor_username"`
	ActorDisplayName *string    `json:"actor_display_name"`
	ActorAvatarURL   *string    `json:"actor_avatar_url"`
	Type             string     `json:"type"`
	EntityID         *uuid.UUID `json:"entity_id"`
	EntityType       *string    `json:"entity_type"`
	Title            string     `json:"title"`
	Body             string     `json:"body"`
	ImageURL         *string    `json:"image_url"`
	ActionURL        *string    `json:"action_url"`
	IsRead           bool       `json:"is_read"`
	ReadAt           *time.Time `json:"read_at"`
	CreatedAt        time.Time  `json:"created_at"`
}

type UserNotificationsResponse struct {
	Notifications []NotificationList `json:"notifications"`
	Pagination    PaginationResponse `json:"pagination"`
}

type PaginationRequest struct {
	Page  uint `json:"page"`
	Limit uint `json:"limit"`
}

type PaginationResponse struct {
	TotalCount uint `json:"total_count"`
	Page       uint `json:"page"`
	Limit      uint `json:"limit"`
	TotalPages uint `json:"total_pages"`
}
