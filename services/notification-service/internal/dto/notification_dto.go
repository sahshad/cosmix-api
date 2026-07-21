package dto

import (
	"time"

	"github.com/google/uuid"
)

type NotificationList struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	ActorID          *uuid.UUID
	ActorUsername    *string
	ActorDisplayName *string
	ActorAvatarURL   *string
	Type             string
	EntityID         *uuid.UUID
	EntityType       *string
	Title            string
	Body             string
	ImageURL         *string
	ActionURL        *string
	IsRead           bool
	ReadAt           *time.Time
	CreatedAt        time.Time
}

type UserNotificationsResponse struct {
	Notifications []NotificationList
	Pagination    PaginationResponse
}

type PaginationRequest struct {
	Page  uint
	Limit uint
}

type PaginationResponse struct {
	TotalCount uint
	Page       uint
	Limit      uint
	TotalPages uint
}
