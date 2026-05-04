package services

import (
	"auth-service/internal/models"
	"auth-service/internal/repositories"
)

type AuthSessionServiceInterface interface {
	Create(s *models.AuthSession) (uint, error)
	FindByRefreshTokenHash(refreshTokenHash string) (*models.AuthSession, error)
	Revoke(refreshTokenHash string) error
}

type AuthSessionService struct {
	sessionRepo repositories.AuthSessionRepositoryInterface
}

func NewAuthSessionService(sessionRepo repositories.AuthSessionRepositoryInterface) AuthSessionServiceInterface {
	return &AuthSessionService{sessionRepo: sessionRepo}
}

func (svc *AuthSessionService) Create(s *models.AuthSession) (uint, error) {
	return svc.sessionRepo.Create(s)
}

func (svc *AuthSessionService) FindByRefreshTokenHash(refreshTokenHash string) (*models.AuthSession, error) {
	return svc.sessionRepo.FindByRefreshTokenHash(refreshTokenHash)
}

func (svc *AuthSessionService) Revoke(refreshTokenHash string) error {
	return svc.sessionRepo.Revoke(refreshTokenHash)
}