package repositories

import (
	"notification-service/internal/models"

	"gorm.io/gorm"
)

type NotificationRepositoryInterface interface {
	Create(notification *models.Notification) error
	GetByUserID(userID uint, limit int, offset int) ([]models.Notification, error)
	GetUnreadCount(userID uint) (int64, error)
	MarkAsRead(notificationID uint, userID uint) error
}

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(
	db *gorm.DB,
) NotificationRepositoryInterface {
	return &NotificationRepository{
		db: db,
	}
}

func (repo *NotificationRepository) Create(notification *models.Notification) error {
	return repo.db.Create(notification).Error
}

func (repo *NotificationRepository) GetByUserID(userID uint, limit int, offset int) ([]models.Notification, error) {
	var notifications []models.Notification

	err := repo.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).Error

	return notifications, err
}

func (repo *NotificationRepository) GetUnreadCount(userID uint) (int64, error) {
	var count int64

	err := repo.db.
		Model(&models.Notification{}).
		Where("user_id = ? AND is_read = false", userID).
		Count(&count).Error

	return count, err
}

func (repo *NotificationRepository) MarkAsRead(notificationID uint, userID uint) error {
	return repo.db.
		Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("is_read", true).Error
}
