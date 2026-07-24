package services

import (
	"context"
	"errors"

	"post-service/internal/dto"
	"post-service/internal/models"
	"post-service/internal/repositories"

	"github.com/google/uuid"
)

type CommentService struct {
	repo     *repositories.CommentRepository
	postRepo *repositories.PostRepository
	userRepo *repositories.PostUserRepository
}

func NewCommentService(
	repo *repositories.CommentRepository,
	postRepo *repositories.PostRepository,
	userRepo *repositories.PostUserRepository,
) *CommentService {
	return &CommentService{
		repo:     repo,
		postRepo: postRepo,
		userRepo: userRepo,
	}
}

func (svc *CommentService) CreateComment(ctx context.Context, postID uuid.UUID, userID uuid.UUID, req *dto.CreateCommentRequest) (*models.Comment, error) {
	_, err := svc.postRepo.FindByID(ctx, postID)
	if err != nil {
		return nil, errors.New("post not found")
	}

	if req.ParentCommentID != nil {
		parent, err := svc.repo.FindByID(ctx, *req.ParentCommentID)
		if err != nil {
			return nil, errors.New("parent comment not found")
		}

		if parent.PostID != postID {
			return nil, errors.New("parent comment does not belong to this post")
		}
	}

	comment := &models.Comment{
		PostID:          postID,
		UserID:          userID,
		ParentCommentID: req.ParentCommentID,
		Content:         req.Content,
	}

	if err := svc.repo.Create(ctx, comment); err != nil {
		return nil, err
	}

	if err := svc.postRepo.IncrementCommentsCount(ctx, postID, 1); err != nil {
		return nil, err
	}

	if req.ParentCommentID != nil {
		if err := svc.repo.IncrementRepliesCount(ctx, *req.ParentCommentID, 1); err != nil {
			return nil, err
		}
	}

	postUser, err := svc.userRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	comment.User = *postUser

	return comment, nil
}

func (svc *CommentService) GetCommentsByPostID(ctx context.Context, postID uuid.UUID, viewerID uuid.UUID, params *dto.PaginationRequest) (*dto.CommentListResponse, error) {
	return svc.repo.GetCommentsByPostID(ctx, postID, viewerID, params)
}

func (svc *CommentService) GetReplies(ctx context.Context, parentCommentID uuid.UUID, viewerID uuid.UUID, params *dto.PaginationRequest) (*dto.CommentListResponse, error) {
	return svc.repo.GetReplies(ctx, parentCommentID, viewerID, params)
}

func (svc *CommentService) UpdateComment(ctx context.Context, id uuid.UUID, userID uuid.UUID, req *dto.UpdateCommentRequest) (*models.Comment, error) {
	comment, err := svc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if comment.UserID != userID {
		return nil, errors.New("unauthorized to update this comment")
	}

	comment.Content = req.Content

	if err := svc.repo.Update(ctx, comment); err != nil {
		return nil, err
	}
	return comment, nil
}

func (svc *CommentService) DeleteComment(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	comment, err := svc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if comment.UserID != userID {
		return errors.New("unauthorized to delete this comment")
	}

	descendantCount, err := svc.repo.CountDescendants(ctx, id)
	if err != nil {
		return err
	}

	if err := svc.repo.DeleteByID(ctx, id); err != nil {
		return err
	}

	if err := svc.postRepo.IncrementCommentsCount(ctx, comment.PostID, -(1 + descendantCount)); err != nil {
		return err
	}

	if comment.ParentCommentID != nil {
		if err := svc.repo.IncrementRepliesCount(ctx, *comment.ParentCommentID, -1); err != nil {
			return err
		}
	}

	return nil
}
