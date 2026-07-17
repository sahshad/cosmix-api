package repositories

import (
	"context"

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
	var notifications []dto.NotificationList

	err := repo.db.WithContext(ctx).
		Table("notifications").
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

	return &dto.UserNotificationsResponse{Notifications: notifications, Pagination: dto.PaginationResponse{
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
