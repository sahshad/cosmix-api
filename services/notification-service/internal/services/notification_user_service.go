package services

import (
	"context"

	"notification-service/internal/models"
	"notification-service/internal/repositories"

	authEvents "cosmix/shared/events/auth"
	userEvents "cosmix/shared/events/user"
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

func (svc *NotificationUserService) HandleAuthUserRegistered(ctx context.Context, event authEvents.AuthUserRegistered) error {
	NotificationUser := &models.NotificationUser{
		UserID:      event.AuthUserID,
		Username:    event.Username,
		DisplayName: event.DisplayName,
		CreatedAt:   event.CreatedAt,
	}

	if err := svc.repo.Create(ctx, NotificationUser); err != nil {
		return err
	}

	return nil
}

func (svc *NotificationUserService) HandleUserUpdated(ctx context.Context, event userEvents.UserUpdated) error {
	user, err := svc.repo.FindByID(ctx, event.UserID)
	if err != nil {
		return err
	}

	user.DisplayName = event.DisplayName
	user.Username = event.Username
	user.AvatarURL = event.AvatarURL
	user.UpdatedAt = &event.UpdatedAt

	if err := svc.repo.Update(ctx, user); err != nil {
		return err
	}
	return nil
}

func (svc *NotificationUserService) HandleUserDeleted(ctx context.Context, event userEvents.UserDeleted) error {
	return svc.repo.DeleteByID(ctx, event.UserID)
}

func (svc *NotificationUserService) FindByUserID(ctx context.Context, userID uint) (*models.NotificationUser, error) {
	return svc.repo.FindByUserID(ctx, userID)
}

func (svc *NotificationUserService) Create(ctx context.Context, notificationUser *models.NotificationUser) error {
	return svc.repo.Create(ctx, notificationUser)
}
