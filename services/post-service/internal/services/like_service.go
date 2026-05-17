package services

import (
	"context"
	"errors"
	"post-service/internal/models"
	"post-service/internal/repositories"
)

type LikeServiceInterface interface {
	LikePost(ctx context.Context, postID uint, userID uint) error
	UnlikePost(ctx context.Context, postID uint, userID uint) error
	GetLikesCount(ctx context.Context, postID uint) (int64, error)
}

type LikeService struct {
	repo     repositories.LikeRepositoryInterface
	postRepo repositories.PostRepositoryInterface
}

func NewLikeService(repo repositories.LikeRepositoryInterface, postRepo repositories.PostRepositoryInterface) LikeServiceInterface {
	return &LikeService{repo: repo, postRepo: postRepo}
}

func (svc *LikeService) LikePost(ctx context.Context, postID uint, userID uint) error {
	// Verify post exists
	_, err := svc.postRepo.FindByID(ctx, postID)
	if err != nil {
		return errors.New("post not found")
	}

	exists, err := svc.repo.Exists(ctx, postID, userID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("already liked")
	}

	like := &models.Like{
		PostID: postID,
		UserID: userID,
	}
	return svc.repo.Create(ctx, like)
}

func (svc *LikeService) UnlikePost(ctx context.Context, postID uint, userID uint) error {
	return svc.repo.Delete(ctx, postID, userID)
}

func (svc *LikeService) GetLikesCount(ctx context.Context, postID uint) (int64, error) {
	return svc.repo.CountByPostID(ctx, postID)
}
