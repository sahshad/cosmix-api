package services

import (
	"context"
	"notification-service/internal/models"
	"notification-service/internal/repositories"
)

// type NotificationPreferenceServiceInterface interface {
// 	CreateDefault(userID uint) error
// 	GetByUserID(userID uint) (*models.NotificationPreference, error)
// 	Update(preference *models.NotificationPreference) error
// }

type NotificationPreferenceService struct {
	repo *repositories.NotificationPreferenceRepository
}

func NewNotificationPreferenceService(
	repo *repositories.NotificationPreferenceRepository,
) *NotificationPreferenceService {
	return &NotificationPreferenceService{
		repo: repo,
	}
}

func (svc *NotificationPreferenceService) CreateDefault(ctx context.Context, userID uint) error {
	preference := &models.NotificationPreference{
		UserID:                 userID,
		EmailEnabled:           true,
		PushEnabled:            true,
		InternalEnabled:        true,
		LikeNotifications:      true,
		CommentNotifications:   true,
		FollowNotifications:    true,
		MentionNotifications:   true,
		MessageNotifications:   true,
		MarketingEmailsEnabled: false,
	}

	return svc.repo.Create(ctx, preference)
}

func (svc *NotificationPreferenceService) GetByUserID(userID uint) (*models.NotificationPreference, error) {
	return svc.repo.GetByUserID(userID)
}

func (svc *NotificationPreferenceService) Update(ctx context.Context, preference *models.NotificationPreference) error {
	return svc.repo.Update(ctx, preference)
}
