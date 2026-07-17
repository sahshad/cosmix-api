package repositories

import (
	"context"

	"post-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// type PostUserRepositoryInterface interface {
// 	Create(ctx context.Context, postUser *models.PostUser) (uuid.UUID, error)
// 	FindByUserID(ctx context.Context, userID uuid.UUID) (*models.PostUser, error)
// 	Update(ctx context.Context, postUser *models.PostUser) error
// 	Delete(ctx context.Context, userID uuid.UUID) error
// }

type PostUserRepository struct {
	*BaseRepository[models.PostUser]
}

func NewPostUserRepository(
	db *gorm.DB,
) *PostUserRepository {
	return &PostUserRepository{
		NewBaseRepository[models.PostUser](db),
	}
}

func (repo *PostUserRepository) Create(ctx context.Context, postUser *models.PostUser) (uuid.UUID, error) {
	result := repo.db.
		WithContext(ctx).
		Create(postUser)

	if result.Error != nil {
		return uuid.Nil, result.Error
	}
	return postUser.UserID, nil
}

func (repo *PostUserRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*models.PostUser, error) {
	var postUser models.PostUser
	if err := repo.db.
		WithContext(ctx).
		Where("user_id = ?", userID).
		First(&postUser).Error; err != nil {
		return nil, err
	}
	return &postUser, nil
}

// func (repo *PostUserRepository) Update(ctx context.Context, postUser *models.PostUser) error {
// 	return repo.db.
// 		WithContext(ctx).
// 		Save(postUser).Error
// }

func (repo *PostUserRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	return repo.db.
		WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&models.PostUser{}).Error
}
