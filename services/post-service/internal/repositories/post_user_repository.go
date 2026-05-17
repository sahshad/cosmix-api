package repositories

import (
	"context"

	"post-service/internal/models"

	"gorm.io/gorm"
)

type PostUserRepositoryInterface interface {
	Create(ctx context.Context, postUser *models.PostUser) (uint, error)
	FindByUserID(ctx context.Context, userID uint) (*models.PostUser, error)
	Update(ctx context.Context, postUser *models.PostUser) error
	Delete(ctx context.Context, userID uint) error
}

type PostUserRepository struct {
	db *gorm.DB
}

func NewPostUserRepository(db *gorm.DB) PostUserRepositoryInterface {
	return &PostUserRepository{db: db}
}

func (repo *PostUserRepository) Create(ctx context.Context, postUser *models.PostUser) (uint, error) {
	result := repo.db.
		WithContext(ctx).
		Create(postUser)

	if result.Error != nil {
		return 0, result.Error
	}
	return postUser.UserID, nil
}

func (repo *PostUserRepository) FindByUserID(ctx context.Context, userID uint) (*models.PostUser, error) {
	var postUser models.PostUser
	if err := repo.db.
		WithContext(ctx).
		Where("user_id = ?", userID).
		First(&postUser).Error; err != nil {
		return nil, err
	}
	return &postUser, nil
}

func (repo *PostUserRepository) Update(ctx context.Context, postUser *models.PostUser) error {
	return repo.db.
		WithContext(ctx).
		Save(postUser).Error
}

func (repo *PostUserRepository) Delete(ctx context.Context, userID uint) error {
	return repo.db.
		WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&models.PostUser{}).Error
}
