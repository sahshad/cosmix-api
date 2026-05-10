package services

import (
	authEvents "cosmix/shared/events/auth"
	"notification-service/internal/models"
)

type EventServiceInterface interface {
	HandleUserRegistered(event authEvents.AuthUserRegistered) error
}

type EventService struct {
	EmailLogSvc               EmailLogServiceInterface
	NotificationSvc           NotificationServiceInterface
	NotificationPreferenceSvc NotificationPreferenceServiceInterface
}

func NewEventService(
	emailLogSvc EmailLogServiceInterface,
	notificationSvc NotificationServiceInterface,
	notificationPreferenceSvc NotificationPreferenceServiceInterface,
) EventServiceInterface {
	return &EventService{
		EmailLogSvc:               emailLogSvc,
		NotificationSvc:           notificationSvc,
		NotificationPreferenceSvc: notificationPreferenceSvc,
	}
}

func (svc *EventService) HandleUserRegistered(event authEvents.AuthUserRegistered) error {
	if err := svc.NotificationPreferenceSvc.CreateDefault(
		event.AuthUserID,
	); err != nil {
		return err
	}

	notification := &models.Notification{
		UserID: event.AuthUserID,
		Type:   "system.welcome",
		Title:  "Welcome to Cosmix",
		Body:   "Your account has been created successfully.",
	}

	if err := svc.NotificationSvc.Create(notification); err != nil {
		return err
	}

	if err := svc.EmailLogSvc.SendWelcomeEmail(
		event.AuthUserID,
		event.Email,
	); err != nil {
		return err
	}

	return nil
}
