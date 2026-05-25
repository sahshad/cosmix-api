package services

import (
	"context"
	"log"
	"notification-service/internal/models"
	"notification-service/internal/repositories"
	"time"
)

// type EmailLogServiceInterface interface {
// 	SendWelcomeEmail(userID uint, email string) error
// 	SendForgotPasswordEmail(userID uint, email string) error
// }

type EmailLogService struct {
	repo *repositories.EmailLogRepository
}

func NewEmailLogService(
	repo *repositories.EmailLogRepository,
) *EmailLogService {
	return &EmailLogService{
		repo: repo,
	}
}

func (svc *EmailLogService) SendWelcomeEmail(ctx context.Context, userID uint, email string) error {
	now := time.Now()

	emailLog := &models.EmailLog{
		UserID:    &userID,
		Recipient: email,
		Type:      "welcome_email",
		Subject:   "Welcome to Cosmix",
		Template:  "welcome.html",
		Status:    "sent",
		SentAt:    &now,
	}

	log.Println("sending welcome email to:", email)

	return svc.repo.Create(ctx, emailLog)
}

func (svc *EmailLogService) SendForgotPasswordEmail(ctx context.Context, userID uint, email string) error {
	now := time.Now()

	emailLog := &models.EmailLog{
		UserID:    &userID,
		Recipient: email,
		Type:      "forgot_password",
		Subject:   "Reset Your Password",
		Template:  "forgot_password.html",
		Status:    "sent",
		SentAt:    &now,
	}

	log.Println("sending forgot password email to:", email)

	return svc.repo.Create(ctx, emailLog)
}
