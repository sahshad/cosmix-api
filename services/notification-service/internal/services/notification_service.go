package services

import (
	"notification-service/internal/models"
	"notification-service/internal/repositories"
)

type NotificationServiceInterface interface {
	Create(notification *models.Notification) error
	GetByUserID(userID uint, limit int, offset int) ([]models.Notification, error)
	GetUnreadCount(userID uint) (int64, error)
	MarkAsRead(notificationID uint, userID uint) error
}

type NotificationService struct {
	repository repositories.NotificationRepositoryInterface
}

func NewNotificationService(
	repository repositories.NotificationRepositoryInterface,
) NotificationServiceInterface {
	return &NotificationService{
		repository: repository,
	}
}

func (svc *NotificationService) Create(notification *models.Notification) error {
	return svc.repository.Create(notification)
}

func (svc *NotificationService) GetByUserID(userID uint, limit int, offset int) ([]models.Notification, error) {
	return svc.repository.GetByUserID(userID, limit, offset)
}

func (svc *NotificationService) GetUnreadCount(userID uint) (int64, error) {
	return svc.repository.GetUnreadCount(userID)
}

func (s *NotificationService) MarkAsRead(notificationID uint, userID uint) error {
	return s.repository.MarkAsRead(notificationID, userID)
}
