package services

import (
	"context"

	"communication-service/internal/models"
	"communication-service/internal/repositories"

	authEvents "cosmix/shared/events/auth"
	userEvents "cosmix/shared/events/user"

	"github.com/google/uuid"
)

type CommunicationUserService struct {
	repo *repositories.CommunicationUserRepository
}

func NewCommunicationUserService(
	repo *repositories.CommunicationUserRepository,
) *CommunicationUserService {
	return &CommunicationUserService{
		repo: repo,
	}
}

func (svc *CommunicationUserService) HandleAuthUserEmailVerificationCompleted(ctx context.Context, event authEvents.AuthUserEmailVerificationCompleted) error {
	user := &models.CommunicationUser{
		UserID:      event.AuthUserID,
		Username:    event.Username,
		DisplayName: event.DisplayName,
		IsVerified:  true,
		CreatedAt:   event.CreatedAt,
	}

	return svc.repo.Create(ctx, user)
}

func (svc *CommunicationUserService) HandleUserUpdated(ctx context.Context, event userEvents.UserUpdated) error {
	user := &models.CommunicationUser{
		UserID:      event.UserID,
		Username:    event.Username,
		DisplayName: event.DisplayName,
		AvatarURL:   event.AvatarURL,
		IsVerified:  event.IsVerified,
		CreatedAt:   event.UpdatedAt,
		UpdatedAt:   &event.UpdatedAt,
	}

	return svc.repo.Upsert(ctx, user)
}

func (svc *CommunicationUserService) FindByUserID(ctx context.Context, userID uuid.UUID) (*models.CommunicationUser, error) {
	return svc.repo.FindByUserID(ctx, userID)
}
