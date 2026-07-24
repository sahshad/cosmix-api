package services

import (
	"context"
	"errors"

	"post-service/internal/models"
	"post-service/internal/repositories"

	"github.com/google/uuid"
)

type CommentLikeService struct {
	repo        *repositories.CommentLikeRepository
	commentRepo *repositories.CommentRepository
}

func NewCommentLikeService(
	repo *repositories.CommentLikeRepository,
	commentRepo *repositories.CommentRepository,
) *CommentLikeService {
	return &CommentLikeService{
		repo:        repo,
		commentRepo: commentRepo,
	}
}

func (svc *CommentLikeService) LikeComment(ctx context.Context, commentID uuid.UUID, userID uuid.UUID) error {
	_, err := svc.commentRepo.FindByID(ctx, commentID)
	if err != nil {
		return errors.New("comment not found")
	}

	exists, err := svc.repo.Exists(ctx, commentID, userID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("already liked")
	}

	like := &models.CommentLike{
		CommentID: commentID,
		UserID:    userID,
	}
	if err := svc.repo.Create(ctx, like); err != nil {
		return err
	}

	return svc.commentRepo.IncrementLikesCount(ctx, commentID, 1)
}

func (svc *CommentLikeService) UnlikeComment(ctx context.Context, commentID uuid.UUID, userID uuid.UUID) error {
	rows, err := svc.repo.Delete(ctx, commentID, userID)
	if err != nil {
		return err
	}
	if rows == 0 {
		return nil
	}

	return svc.commentRepo.IncrementLikesCount(ctx, commentID, -1)
}
