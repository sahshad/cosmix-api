package repositories

import (
	"context"

	"user-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	*BaseRepository[models.User]
}

func NewUserRepository(
	db *gorm.DB,
) *UserRepository {
	return &UserRepository{
		NewBaseRepository[models.User](db),
	}
}

func (repo *UserRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	var profile models.User
	if err := repo.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (repo *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var profile models.User
	if err := repo.db.WithContext(ctx).First(&profile, id).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (repo *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var profile models.User
	if err := repo.db.WithContext(ctx).Where("username = ?", username).First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}
