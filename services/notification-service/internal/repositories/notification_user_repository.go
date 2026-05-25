package repositories

import (
	"context"
	"gorm.io/gorm"
	"notification-service/internal/models"
)

type NotificationUserRepository struct {
	*BaseRepository[models.NotificationUser]
}

func NewNotificationUserRepository(
	db *gorm.DB,
) *NotificationUserRepository {
	return &NotificationUserRepository{
		NewBaseRepository[models.NotificationUser](db),
	}
}

func (repo *NotificationUserRepository) FindByUserID(ctx context.Context, userID uint) (*models.NotificationUser, error) {
	var notificationUser models.NotificationUser
	if err := repo.DB().
		WithContext(ctx).
		Where("user_id = ?", userID).
		First(&notificationUser).Error; err != nil {
		return nil, err
	}
	return &notificationUser, nil
}
