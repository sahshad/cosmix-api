package repositories

import (
	"user-service/internal/models"

	"context"
	"gorm.io/gorm"
)

// type UserProfileRepositoryInterface interface {
// 	Create(ctx context.Context, profile *models.User) error
// 	FindByUserID(ctx context.Context, userID uint) (*models.User, error)
// 	FindByID(ctx context.Context, id uint) (*models.User, error)
// 	FindByUsername(ctx context.Context, username string) (*models.User, error)
// 	Update(ctx context.Context, profile *models.User) error
// 	Delete(ctx context.Context, id uint) error
// }

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

// func (repo *UserRepo) Create(ctx context.Context, profile *models.User) error {
// 	return repo.db.WithContext(ctx).Create(profile).Error
// }

func (repo *UserRepository) FindByUserID(ctx context.Context, userID uint) (*models.User, error) {
	var profile models.User
	if err := repo.db.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (repo *UserRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var profile models.User
	if err := repo.db.WithContext(ctx).First(&profile, id).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

// func (repo *UserRepo) Update(ctx context.Context, profile *models.User) error {
// 	return repo.db.WithContext(ctx).Save(profile).Error
// }

func (repo *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var profile models.User
	if err := repo.db.WithContext(ctx).Where("username = ?", username).First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

// func (repo *UserRepo) Delete(ctx context.Context, id uint) error {
// 	return repo.db.WithContext(ctx).Delete(&models.User{}, id).Error
// }
