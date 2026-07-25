package dto

import (
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
	ReadAt           *string
	CreatedAt        string
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
