package services

import (
	"context"

	"notification-service/internal/email"
	"notification-service/internal/models"

	authEvents "cosmix/shared/events/auth"
)

type EventService struct {
	EmailLogSvc               *EmailLogService
	NotificationSvc           *NotificationService
	NotificationPreferenceSvc *NotificationPreferenceService
	NotificationUserSvc       *NotificationUserService
	MailDispatcher            *email.MailDispatcher
}

func NewEventService(
	emailLogSvc *EmailLogService,
	notificationSvc *NotificationService,
	notificationPreferenceSvc *NotificationPreferenceService,
	notificationUserSvc *NotificationUserService,
	mailDispatcher *email.MailDispatcher,
) *EventService {
	return &EventService{
		EmailLogSvc:               emailLogSvc,
		NotificationSvc:           notificationSvc,
		NotificationPreferenceSvc: notificationPreferenceSvc,
		NotificationUserSvc:       notificationUserSvc,
		MailDispatcher:            mailDispatcher,
	}
}

func (svc *EventService) HandleUserEmailVerificationCompleted(ctx context.Context, event authEvents.AuthUserEmailVerificationCompleted) error {
	notificationUser := &models.NotificationUser{
		UserID:      event.AuthUserID,
		Username:    event.Username,
		DisplayName: event.DisplayName,
		CreatedAt:   event.CreatedAt,
	}
	if err := svc.NotificationUserSvc.Create(ctx, notificationUser); err != nil {
		return err
	}

	if err := svc.NotificationPreferenceSvc.CreateDefault(ctx, event.AuthUserID); err != nil {
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

	// if err := svc.EmailLogSvc.SendWelcomeEmail(ctx, event.AuthUserID, event.Email); err != nil {
	// 	return err
	// }

	appLink := "http://localhost:3000"

	templateData := map[string]any{
		"Title":       "Welcome to Cosmix",
		"DisplayName": event.DisplayName,
		"AppLink":     appLink,
	}

	mailDto := email.SendEmailDTO{
		To:           event.Email,
		Subject:      "Welcome to Cosmix",
		TemplateName: email.TemplateAuthWelcome,
		TemplateData: templateData,
	}

	err := svc.MailDispatcher.Send(mailDto)
	if err != nil {
		return err
	}

	return nil
}
