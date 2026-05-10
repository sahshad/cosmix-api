package services

import (
	"context"
	"notification-service/internal/dto"
	"notification-service/internal/models"
	"notification-service/internal/repositories"
)

type NotificationServiceInterface interface {
	Create(notification *models.Notification) error
	GetUserNotifications(ctx context.Context, userID uint, paginationRequest dto.PaginationRequest) (*dto.UserNotificationsResponse, error)
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

func (svc *NotificationService) GetUserNotifications(ctx context.Context, userID uint, paginationRequest dto.PaginationRequest) (*dto.UserNotificationsResponse, error) {
	return svc.repository.GetUserNotifications(ctx, userID, paginationRequest)
}

func (svc *NotificationService) GetUnreadCount(userID uint) (int64, error) {
	return svc.repository.GetUnreadCount(userID)
}

func (s *NotificationService) MarkAsRead(notificationID uint, userID uint) error {
	return s.repository.MarkAsRead(notificationID, userID)
}
