package services

import (
	"context"

	"notification-service/internal/models"
	"notification-service/internal/repositories"

	userEvents "cosmix/shared/events/user"

	"github.com/google/uuid"
)

type NotificationUserService struct {
	repo *repositories.NotificationUserRepository
}

func NewNotificationUserService(
	repo *repositories.NotificationUserRepository,
) *NotificationUserService {
	return &NotificationUserService{
		repo: repo,
	}
}

func (svc *NotificationUserService) HandleUserUpdated(ctx context.Context, event userEvents.UserUpdated) error {
	user, err := svc.repo.FindByUserID(ctx, event.UserID)
	if err != nil {
		return err
	}

	user.DisplayName = event.DisplayName
	user.Username = event.Username
	user.AvatarURL = event.AvatarURL
	user.IsVerified = event.IsVerified
	user.UpdatedAt = &event.UpdatedAt

	if err := svc.repo.Update(ctx, user); err != nil {
		return err
	}
	return nil
}

func (svc *NotificationUserService) HandleUserDeleted(ctx context.Context, event userEvents.UserDeleted) error {
	return svc.repo.DeleteByID(ctx, event.UserID)
}

func (svc *NotificationUserService) FindByUserID(ctx context.Context, userID uuid.UUID) (*models.NotificationUser, error) {
	return svc.repo.FindByUserID(ctx, userID)
}

func (svc *NotificationUserService) Create(ctx context.Context, notificationUser *models.NotificationUser) error {
	return svc.repo.Create(ctx, notificationUser)
}
