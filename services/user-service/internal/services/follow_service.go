package services

import (
	"context"
	"errors"
	"user-service/internal/models"
	"user-service/internal/repositories"
)

// type FollowServiceInterface interface {
// 	Follow(ctx context.Context, followerID uint, followingID uint) error
// 	Unfollow(ctx context.Context, followerID uint, followingID uint) error
// 	GetFollowers(ctx context.Context, userID uint) ([]uint, error)
// 	GetFollowing(ctx context.Context, userID uint) ([]uint, error)
// }

type FollowService struct {
	repo *repositories.FollowRepository
}

func NewFollowService(
	repo *repositories.FollowRepository,
	) *FollowService {
	return &FollowService{
		repo: repo,
	}
}

func (svc *FollowService) Follow(ctx context.Context, followerID, followingID uint) error {
	if followerID == followingID {
		return errors.New("cannot follow yourself")
	}

	isFollowing, err := svc.repo.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return err
	}
	if isFollowing {
		return errors.New("already following")
	}

	follow := &models.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}
	return svc.repo.Create(ctx, follow)
}

func (svc *FollowService) Unfollow(ctx context.Context, followerID, followingID uint) error {
	return svc.repo.Delete(ctx, followerID, followingID)
}

func (svc *FollowService) GetFollowers(ctx context.Context, userID uint) ([]uint, error) {
	return svc.repo.GetFollowers(ctx, userID)
}

func (svc *FollowService) GetFollowing(ctx context.Context, userID uint) ([]uint, error) {
	return svc.repo.GetFollowing(ctx, userID)
}
