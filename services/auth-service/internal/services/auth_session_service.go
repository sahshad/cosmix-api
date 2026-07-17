package services

import (
	"context"

	"auth-service/internal/models"
	"auth-service/internal/repositories"

	"github.com/google/uuid"
)

type AuthSessionService struct {
	repo *repositories.AuthSessionRepository
}

func NewAuthSessionService(
	repo *repositories.AuthSessionRepository,
) *AuthSessionService {
	return &AuthSessionService{
		repo: repo,
	}
}

func (svc *AuthSessionService) Create(ctx context.Context, authSession *models.AuthSession) (uuid.UUID, error) {
	err := svc.repo.Create(ctx, authSession)
	if err != nil {
		return uuid.Nil, err
	}

	return authSession.ID, nil
}

func (svc *AuthSessionService) FindByRefreshTokenHash(refreshTokenHash string) (*models.AuthSession, error) {
	return svc.repo.FindByRefreshTokenHash(refreshTokenHash)
}

func (svc *AuthSessionService) Revoke(refreshTokenHash string) error {
	return svc.repo.Revoke(refreshTokenHash)
}
