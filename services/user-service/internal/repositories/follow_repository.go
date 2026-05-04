package repositories

import (
	"context"
	"user-service/internal/models"

	"gorm.io/gorm"
)

type FollowRepositoryInterface interface {
	Create(ctx context.Context, follow *models.Follow) error
	Delete(ctx context.Context, followerID, followingID uint) error
	IsFollowing(ctx context.Context, followerID, followingID uint) (bool, error)
	GetFollowers(ctx context.Context, userID uint) ([]uint, error)
	GetFollowing(ctx context.Context, userID uint) ([]uint, error)
	GetFollowerCount(ctx context.Context, userID uint) (int64, error)
	GetFollowingCount(ctx context.Context, userID uint) (int64, error)
}

type FollowRepository struct {
	db *gorm.DB
}

func NewFollowRepository(db *gorm.DB) *FollowRepository {
	return &FollowRepository{db: db}
}

func (repo *FollowRepository) Create(ctx context.Context, follow *models.Follow) error {
	return repo.db.WithContext(ctx).Create(follow).Error
}

func (repo *FollowRepository) Delete(ctx context.Context, followerID, followingID uint) error {
	return repo.db.WithContext(ctx).Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Delete(&models.Follow{}).Error
}

func (repo *FollowRepository) IsFollowing(ctx context.Context, followerID, followingID uint) (bool, error) {
	var count int64
	err := repo.db.WithContext(ctx).Model(&models.Follow{}).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Count(&count).Error
	return count > 0, err
}

func (repo *FollowRepository) GetFollowers(ctx context.Context, userID uint) ([]uint, error) {
	var followers []uint
	err := repo.db.WithContext(ctx).Model(&models.Follow{}).
		Where("following_id = ?", userID).
		Pluck("follower_id", &followers).Error
	return followers, err
}

func (repo *FollowRepository) GetFollowing(ctx context.Context, userID uint) ([]uint, error) {
	var following []uint
	err := repo.db.WithContext(ctx).Model(&models.Follow{}).
		Where("follower_id = ?", userID).
		Pluck("following_id", &following).Error
	return following, err
}

func (repo *FollowRepository) GetFollowerCount(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := repo.db.WithContext(ctx).Model(&models.Follow{}).
		Where("following_id = ?", userID).
		Count(&count).Error
	return count, err
}

func (repo *FollowRepository) GetFollowingCount(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := repo.db.WithContext(ctx).Model(&models.Follow{}).
		Where("follower_id = ?", userID).
		Count(&count).Error
	return count, err
}
