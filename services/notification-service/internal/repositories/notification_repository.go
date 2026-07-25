package repositories

import (
	"context"
	"time"

	"notification-service/internal/dto"
	"notification-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	*BaseRepository[models.Notification]
}

func NewNotificationRepository(
	db *gorm.DB,
) *NotificationRepository {
	return &NotificationRepository{
		NewBaseRepository[models.Notification](db),
	}
}

func (repo *NotificationRepository) GetUserNotifications(ctx context.Context, userID uuid.UUID, params dto.PaginationRequest) (*dto.UserNotificationsResponse, error) {
	var notifications []models.Notification

	err := repo.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(int(params.Limit)).
		Offset(int((params.Page - 1) * params.Limit)).
		Find(&notifications).Error

	if err != nil {
		return nil, err
	}

	totalCount := uint(len(notifications))
	totalPages := (uint(len(notifications)) + params.Limit - 1) / params.Limit

	notificationList := make([]dto.NotificationList, 0, len(notifications))
	for _, n := range notifications {
		notificationList = append(notificationList, dto.NotificationList{
			ID:               n.ID,
			UserID:           n.UserID,
			ActorID:          n.ActorID,
			ActorUsername:    n.ActorUsername,
			ActorDisplayName: n.ActorDisplayName,
			ActorAvatarURL:   n.ActorAvatarURL,
			Type:             n.Type,
			EntityID:         n.EntityID,
			EntityType:       n.EntityType,
			Title:            n.Title,
			Body:             n.Body,
			ImageURL:         n.ImageURL,
			ActionURL:        n.ActionURL,
			IsRead:           n.IsRead,
			ReadAt:           formatTimePtr(n.ReadAt),
			CreatedAt:        n.CreatedAt.Format(time.RFC3339),
		})
	}

	return &dto.UserNotificationsResponse{Notifications: notificationList, Pagination: dto.PaginationResponse{
		Page:       params.Page,
		Limit:      params.Limit,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}}, nil
}

func (repo *NotificationRepository) GetUnreadCount(userID uuid.UUID) (int64, error) {
	var count int64

	err := repo.db.
		Model(&models.Notification{}).
		Where("user_id = ? AND is_read = false", userID).
		Count(&count).Error

	return count, err
}

func (repo *NotificationRepository) MarkAsRead(notificationID uuid.UUID, userID uuid.UUID) error {
	return repo.db.
		Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("is_read", true).Error
}

func formatTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format(time.RFC3339)
	return &s
}
