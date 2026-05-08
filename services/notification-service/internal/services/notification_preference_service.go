package services

import (
	"notification-service/internal/models"
	"notification-service/internal/repositories"
)

type NotificationPreferenceServiceInterface interface {
	CreateDefault(userID uint) error
	GetByUserID(userID uint) (*models.NotificationPreference, error)
	Update(preference *models.NotificationPreference) error
}

type NotificationPreferenceService struct {
	repository repositories.NotificationPreferenceRepositoryInterface
}

func NewNotificationPreferenceService(
	repository repositories.NotificationPreferenceRepositoryInterface,
) NotificationPreferenceServiceInterface {
	return &NotificationPreferenceService{
		repository: repository,
	}
}

func (svc *NotificationPreferenceService) CreateDefault(userID uint) error {
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

	return svc.repository.Create(preference)
}

func (svc *NotificationPreferenceService) GetByUserID(userID uint) (*models.NotificationPreference, error) {
	return svc.repository.GetByUserID(userID)
}

func (svc *NotificationPreferenceService) Update(preference *models.NotificationPreference) error {
	return svc.repository.Update(preference)
}
