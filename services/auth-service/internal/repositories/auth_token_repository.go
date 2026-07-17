package repositories

import (
	"auth-service/internal/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthTokenRepository struct {
	*BaseRepository[models.AuthToken]
}

func NewAuthTokenRepository(
	db *gorm.DB,
) *AuthTokenRepository {
	return &AuthTokenRepository{
		NewBaseRepository[models.AuthToken](db),
	}
}

func (repo *AuthTokenRepository) FindByToken(ctx context.Context, token string) (*models.AuthToken, error) {
	var authToken models.AuthToken
	if err := repo.db.WithContext(ctx).Where("token = ?", token).First(&authToken).Error; err != nil {
		return nil, err
	}
	return &authToken, nil
}

func (repo *AuthTokenRepository) FindByAuthUserID(ctx context.Context, authUserID uuid.UUID) (*models.AuthToken, error) {
	var authToken models.AuthToken
	if err := repo.db.WithContext(ctx).Where("auth_user_id = ?", authUserID).Preload("AuthUser").First(&authToken).Error; err != nil {
		return nil, err
	}
	return &authToken, nil
}
