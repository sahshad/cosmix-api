package repositories

import (
	"auth-service/internal/models"
	"context"

	"gorm.io/gorm"
)

type EmailVerificationTokenRepository struct {
	*BaseRepository[models.EmailVerificationToken]
}

func NewEmailVerificationTokenRepository(
	db *gorm.DB,
) *EmailVerificationTokenRepository {
	return &EmailVerificationTokenRepository{
		NewBaseRepository[models.EmailVerificationToken](db),
	}
}

func (repo *EmailVerificationTokenRepository) FindByToken(ctx context.Context, token string) (*models.EmailVerificationToken, error) {
	var emailVerificationToken models.EmailVerificationToken
	if err := repo.db.WithContext(ctx).Where("token = ?", token).First(&emailVerificationToken).Error; err != nil {
		return nil, err
	}
	return &emailVerificationToken, nil
}

func (repo *EmailVerificationTokenRepository) FindByAuthUserID(ctx context.Context, authUserID uint) (*models.EmailVerificationToken, error) {
	var emailVerificationToken models.EmailVerificationToken
	if err := repo.db.WithContext(ctx).Where("auth_user_id = ?", authUserID).Preload("AuthUser").First(&emailVerificationToken).Error; err != nil {
		return nil, err
	}
	return &emailVerificationToken, nil
}
