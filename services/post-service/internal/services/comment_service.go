package services

import (
	"context"
	"errors"

	"post-service/internal/dto"
	"post-service/internal/models"
	"post-service/internal/repositories"
)

// type CommentServiceInterface interface {
// 	CreateComment(ctx context.Context, postID uint, userID uint, req *dto.CreateCommentRequest) (*models.Comment, error)
// 	GetCommentsByPostID(ctx context.Context, postID uint, params *dto.PaginationRequest) (*dto.CommentListResponse, error)
// 	UpdateComment(ctx context.Context, id uint, userID uint, req *dto.UpdateCommentRequest) (*models.Comment, error)
// 	DeleteComment(ctx context.Context, id uint, userID uint) error
// }

type CommentService struct {
	repo     *repositories.CommentRepository
	postRepo *repositories.PostRepository
}

func NewCommentService(
	repo *repositories.CommentRepository,
	postRepo *repositories.PostRepository,
) *CommentService {
	return &CommentService{
		repo:     repo,
		postRepo: postRepo,
	}
}

func (svc *CommentService) CreateComment(ctx context.Context, postID uint, userID uint, req *dto.CreateCommentRequest) (*models.Comment, error) {
	_, err := svc.postRepo.FindByID(ctx, postID)
	if err != nil {
		return nil, errors.New("post not found")
	}

	comment := &models.Comment{
		PostID:  postID,
		UserID:  userID,
		Content: req.Content,
	}

	if err := svc.repo.Create(ctx, comment); err != nil {
		return nil, err
	}
	return comment, nil
}

func (svc *CommentService) GetCommentsByPostID(ctx context.Context, postID uint, params *dto.PaginationRequest) (*dto.CommentListResponse, error) {
	return svc.repo.GetCommentsByPostID(ctx, postID, params)
}

func (svc *CommentService) UpdateComment(ctx context.Context, id uint, userID uint, req *dto.UpdateCommentRequest) (*models.Comment, error) {
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

func (svc *CommentService) DeleteComment(ctx context.Context, id uint, userID uint) error {
	comment, err := svc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if comment.UserID != userID {
		return errors.New("unauthorized to delete this comment")
	}

	return svc.repo.DeleteByID(ctx, id)
}
