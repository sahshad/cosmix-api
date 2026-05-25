package services

import (
	"context"
	authEvents "cosmix/shared/events/auth"
	"notification-service/internal/models"
)

// type EventServiceInterface interface {
// 	HandleUserRegistered(event authEvents.AuthUserRegistered) error
// }

type EventService struct {
	EmailLogSvc               *EmailLogService
	NotificationSvc           *NotificationService
	NotificationPreferenceSvc *NotificationPreferenceService
	NotificationUserSvc *NotificationUserService
}

func NewEventService(
	emailLogSvc *EmailLogService,
	notificationSvc *NotificationService,
	notificationPreferenceSvc *NotificationPreferenceService,
	notificationUserSvc *NotificationUserService,
) *EventService {
	return &EventService{
		EmailLogSvc:               emailLogSvc,
		NotificationSvc:           notificationSvc,
		NotificationPreferenceSvc: notificationPreferenceSvc,
		NotificationUserSvc: notificationUserSvc,
	}
}

func (svc *EventService) HandleUserRegistered(ctx context.Context, event authEvents.AuthUserRegistered) error {
	notificationUser := &models.NotificationUser{
		UserID: event.AuthUserID,
		Username: event.Username,
		DisplayName: event.DisplayName,
		CreatedAt: event.CreatedAt,
	}
	if err := svc.NotificationUserSvc.Create(ctx, notificationUser); err != nil {
		return err
	}

	
	if err := svc.NotificationPreferenceSvc.CreateDefault(
		ctx,
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

	if err := svc.NotificationSvc.Create(ctx, notification); err != nil {
		return err
	}

	if err := svc.EmailLogSvc.SendWelcomeEmail(
		ctx,
		event.AuthUserID,
		event.Email,
	); err != nil {
		return err
	}

	return nil
}
