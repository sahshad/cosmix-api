package repositories

import (
	"context"
	
	"post-service/internal/models"

	"gorm.io/gorm"
)

// type LikeRepositoryInterface interface {
// 	Create(ctx context.Context, like *models.Like) error
// 	Delete(ctx context.Context, postID uint, userID uint) error
// 	CountByPostID(ctx context.Context, postID uint) (int64, error)
// 	Exists(ctx context.Context, postID uint, userID uint) (bool, error)
// }

type LikeRepository struct {
	*BaseRepository[models.Like]
}

func NewLikeRepository(
	db *gorm.DB,
	) *LikeRepository {
	return &LikeRepository{
	 NewBaseRepository[models.Like](db),
	}
}

// func (repo *LikeRepository) Create(ctx context.Context, like *models.Like) error {
// 	return repo.db.WithContext(ctx).Create(like).Error
// }

func (repo *LikeRepository) Delete(ctx context.Context, postID uint, userID uint) error {
	return repo.db.WithContext(ctx).Where("post_id = ? AND user_id = ?", postID, userID).Delete(&models.Like{}).Error
}

func (repo *LikeRepository) CountByPostID(ctx context.Context, postID uint) (int64, error) {
	var count int64
	err := repo.db.WithContext(ctx).Model(&models.Like{}).Where("post_id = ?", postID).Count(&count).Error
	return count, err
}

func (repo *LikeRepository) Exists(ctx context.Context, postID uint, userID uint) (bool, error) {
	var count int64
	err := repo.db.WithContext(ctx).Model(&models.Like{}).Where("post_id = ? AND user_id = ?", postID, userID).Count(&count).Error
	return count > 0, err
}
