package services

import (
	"context"
	"notification-service/internal/dto"
	"notification-service/internal/models"
	"notification-service/internal/repositories"
)

// type NotificationServiceInterface interface {
// 	Create(notification *models.Notification) error
// 	GetUserNotifications(ctx context.Context, userID uint, paginationRequest dto.PaginationRequest) (*dto.UserNotificationsResponse, error)
// 	GetUnreadCount(userID uint) (int64, error)
// 	MarkAsRead(notificationID uint, userID uint) error
// }

type NotificationService struct {
	repo *repositories.NotificationRepository
}

func NewNotificationService(
	repo *repositories.NotificationRepository,
) *NotificationService {
	return &NotificationService{
		repo: repo,
	}
}

func (svc *NotificationService) Create(ctx context.Context, notification *models.Notification) error {
	return svc.repo.Create(ctx, notification)
}

func (svc *NotificationService) GetUserNotifications(ctx context.Context, userID uint, paginationRequest dto.PaginationRequest) (*dto.UserNotificationsResponse, error) {
	return svc.repo.GetUserNotifications(ctx, userID, paginationRequest)
}

func (svc *NotificationService) GetUnreadCount(userID uint) (int64, error) {
	return svc.repo.GetUnreadCount(userID)
}

func (s *NotificationService) MarkAsRead(notificationID uint, userID uint) error {
	return s.repo.MarkAsRead(notificationID, userID)
}
