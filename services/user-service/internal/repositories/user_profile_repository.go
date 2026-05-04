package repositories

import (
	"user-service/internal/models"

	"context"
	"gorm.io/gorm"
)

type UserProfileRepositoryInterface interface {
	Create(ctx context.Context, profile *models.UserProfile) error
	FindByUserID(ctx context.Context, userID uint) (*models.UserProfile, error)
	FindByID(ctx context.Context, id uint) (*models.UserProfile, error)
	FindByUsername(ctx context.Context, username string) (*models.UserProfile, error)
	Update(ctx context.Context, profile *models.UserProfile) error
	Delete(ctx context.Context, id uint) error
}

type UserProfileRepo struct {
	db *gorm.DB
}

func NewUserProfileRepository(db *gorm.DB) *UserProfileRepo {
	return &UserProfileRepo{db: db}
}

func (repo *UserProfileRepo) Create(ctx context.Context, profile *models.UserProfile) error {
	return repo.db.WithContext(ctx).Create(profile).Error
}

func (repo *UserProfileRepo) FindByUserID(ctx context.Context, userID uint) (*models.UserProfile, error) {
	var profile models.UserProfile
	if err := repo.db.Where("auth_user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (repo *UserProfileRepo) FindByID(ctx context.Context, id uint) (*models.UserProfile, error) {
	var profile models.UserProfile
	if err := repo.db.WithContext(ctx).First(&profile, id).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (repo *UserProfileRepo) Update(ctx context.Context, profile *models.UserProfile) error {
	return repo.db.WithContext(ctx).Save(profile).Error
}

func (repo *UserProfileRepo) FindByUsername(ctx context.Context, username string) (*models.UserProfile, error) {
	var profile models.UserProfile
	if err := repo.db.WithContext(ctx).Where("username = ?", username).First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (repo *UserProfileRepo) Delete(ctx context.Context, id uint) error {
	return repo.db.WithContext(ctx).Delete(&models.UserProfile{}, id).Error
}
