package services

import (
	"context"

	"user-service/internal/models"
	"user-service/internal/repositories"

	appErr "cosmix/shared/core/errors"

	"github.com/google/uuid"
)

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

func (svc *FollowService) Follow(ctx context.Context, followerID, followingID uuid.UUID) error {
	if followerID == followingID {
		return appErr.NewBadRequest("followingId", "cannot follow yourself")
	}

	isFollowing, err := svc.repo.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return err
	}
	if isFollowing {
		return appErr.NewConflict("already following")
	}

	follow := &models.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}
	if err := svc.repo.Create(ctx, follow); err != nil {
		return err
	}

	if err := svc.repo.IncrementFollowingCount(ctx, followerID, 1); err != nil {
		return err
	}
	return svc.repo.IncrementFollowersCount(ctx, followingID, 1)
}

func (svc *FollowService) Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error {
	rows, err := svc.repo.Delete(ctx, followerID, followingID)
	if err != nil {
		return err
	}
	if rows == 0 {
		return nil
	}

	if err := svc.repo.IncrementFollowingCount(ctx, followerID, -1); err != nil {
		return err
	}
	return svc.repo.IncrementFollowersCount(ctx, followingID, -1)
}

func (svc *FollowService) GetFollowers(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	return svc.repo.GetFollowers(ctx, userID)
}

func (svc *FollowService) GetFollowing(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	return svc.repo.GetFollowing(ctx, userID)
}
