package services

import (
	"context"
	"fmt"

	"notification-service/internal/dto"
	"notification-service/internal/email"
	"notification-service/internal/models"
	"notification-service/internal/repositories"

	authEvents "cosmix/shared/events/auth"
)

type NotificationService struct {
	repo           *repositories.NotificationRepository
	mailDispatcher *email.MailDispatcher
}

func NewNotificationService(
	repo *repositories.NotificationRepository,
	mailDispatcher *email.MailDispatcher,
) *NotificationService {
	return &NotificationService{
		repo:           repo,
		mailDispatcher: mailDispatcher,
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

func (svc *NotificationService) MarkAsRead(notificationID uint, userID uint) error {
	return svc.repo.MarkAsRead(notificationID, userID)
}

func (svc *NotificationService) HandleEmailVerification(ctx context.Context, event authEvents.AuthUserEmailVerification) error {
	verificationLink := fmt.Sprintf("http://localhost:3000/verify-email?token=%s", event.Token)
	templateData := map[string]any{
		"Title":            "Verify Your Email",
		"DisplayName":     event.DisplayName,
		"VerificationLink": verificationLink,
	}

	mailDto := email.SendEmailDTO{
		To:           event.Email,
		Subject:      "Verify your email",
		TemplateName: email.TemplateAuthVerification,
		TemplateData: templateData,
	}

	err := svc.mailDispatcher.Send(mailDto)
	if err != nil {
		return err
	}

	return nil
}

// func (svc *NotificationService) SendEmailVerification(to string, token string) error {
// 	verificationLink := fmt.Sprintf("http://localhost:3000/verify-email?token=%s", token)
// 	err := svc.mailDispatcher.Send(
// 		email.SendEmailDTO{
// 			To:           to,
// 			Subject:      "Verify your email",
// 			TemplateName: email.TemplateAuthVerification,

// 			TemplateData: map[string]any{
// 				"Title": "Verify Your Email",

// 				"VerificationLink": verificationLink,
// 			},
// 		},
// 	)

// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
