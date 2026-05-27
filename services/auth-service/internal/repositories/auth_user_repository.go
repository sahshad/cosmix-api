package repositories

import (
	"context"

	"auth-service/internal/models"

	"gorm.io/gorm"
)

type AuthUserRepository struct {
	*BaseRepository[models.AuthUser]
}

func NewAuthUserRepository(
	db *gorm.DB,
) *AuthUserRepository {
	return &AuthUserRepository{
		NewBaseRepository[models.AuthUser](db),
	}
}

func (repo *AuthUserRepository) FindByEmail(ctx context.Context, email string) (*models.AuthUser, error) {
	var user models.AuthUser
	if err := repo.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *AuthUserRepository) FindByUsername(ctx context.Context, username string) (*models.AuthUser, error) {
	var user models.AuthUser
	if err := repo.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
