package services

import (
	"context"
	"errors"
	"fmt"

	"post-service/internal/models"
	"post-service/internal/repositories"

	"github.com/google/uuid"
)

// type LikeServiceInterface interface {
// 	LikePost(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error
// 	UnlikePost(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error
// 	GetLikesCount(ctx context.Context, postID uuid.UUID) (int64, error)
// }

type LikeService struct {
	repo     *repositories.LikeRepository
	postRepo *repositories.PostRepository
}

func NewLikeService(
	repo *repositories.LikeRepository,
	postRepo *repositories.PostRepository,
) *LikeService {
	return &LikeService{
		repo:     repo,
		postRepo: postRepo,
	}
}

func (svc *LikeService) LikePost(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error {
	fmt.Println("in the like service")
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
	if err := svc.repo.Create(ctx, like); err != nil {
		return err
	}

	return svc.postRepo.IncrementLikesCount(ctx, postID, 1)
}

func (svc *LikeService) UnlikePost(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error {
	rows, err := svc.repo.Delete(ctx, postID, userID)
	if err != nil {
		return err
	}
	if rows == 0 {
		return nil
	}

	return svc.postRepo.IncrementLikesCount(ctx, postID, -1)
}

func (svc *LikeService) GetLikesCount(ctx context.Context, postID uuid.UUID) (int64, error) {
	return svc.repo.CountByPostID(ctx, postID)
}
