package services

import (
	"log"
	"notification-service/internal/models"
	"notification-service/internal/repositories"
	"time"
)

type EmailLogServiceInterface interface {
	SendWelcomeEmail(userID uint, email string) error
	SendForgotPasswordEmail(userID uint, email string) error
}

type EmailLogService struct {
	repository repositories.EmailLogRepositoryInterface
}

func NewEmailLogService(
	repository repositories.EmailLogRepositoryInterface,
) EmailLogServiceInterface {
	return &EmailLogService{
		repository: repository,
	}
}

func (svc *EmailLogService) SendWelcomeEmail(userID uint, email string) error {
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

	return svc.repository.Create(emailLog)
}

func (svc *EmailLogService) SendForgotPasswordEmail(userID uint, email string) error {
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

	return svc.repository.Create(emailLog)
}
