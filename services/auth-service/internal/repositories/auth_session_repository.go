package repositories

import (
	"auth-service/internal/models"

	"gorm.io/gorm"
)

// type AuthSessionRepositoryInterface interface {
// 	Create(s *models.AuthSession) (uint, error)
// 	FindByRefreshTokenHash(refreshTokenHash string) (*models.AuthSession, error)
// 	Revoke(refreshTokenHash string) error
// }

type AuthSessionRepository struct {
	*BaseRepository[models.AuthSession]
}

func NewAuthSessionRepository(
	db *gorm.DB,
	) *AuthSessionRepository {
	return &AuthSessionRepository{
		NewBaseRepository[models.AuthSession](db),
	}
}

// func (repo *AuthSessionRepo) Create(s *models.AuthSession) (uint, error) {
// 	if err := repo.db.Create(s).Error; err != nil {
// 		return 0, err
// 	}
// 	return s.ID, nil
// }

func (repo *AuthSessionRepository) FindByRefreshTokenHash(refreshTokenHash string) (*models.AuthSession, error) {
	var session models.AuthSession
	if err := repo.db.Where("refresh_token_hash = ?", refreshTokenHash).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (repo *AuthSessionRepository) Revoke(refreshTokenHash string) error {
	return repo.db.Where("refresh_token_hash = ?", refreshTokenHash).Update("revoked", true).Error
}