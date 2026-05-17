package repositories

import (
	"context"
	"post-service/internal/dto"
	"post-service/internal/models"

	"gorm.io/gorm"
)

type CommentRepositoryInterface interface {
	Create(ctx context.Context, comment *models.Comment) error
	GetCommentsByPostID(ctx context.Context, postID uint, params *dto.PaginationRequest) (*dto.CommentListResponse, error)
	FindByID(ctx context.Context, id uint) (*models.Comment, error)
	Update(ctx context.Context, comment *models.Comment) error
	Delete(ctx context.Context, id uint) error
}

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepositoryInterface {
	return &CommentRepository{db: db}
}

func (repo *CommentRepository) Create(ctx context.Context, comment *models.Comment) error {
	return repo.db.WithContext(ctx).Create(comment).Error
}

func (repo *CommentRepository) GetCommentsByPostID(ctx context.Context, postID uint, params *dto.PaginationRequest) (*dto.CommentListResponse, error) {
	var comments []dto.CommentList
	if err := repo.db.WithContext(ctx).
		Table("comments").
		Where("post_id = ?", postID).
		Limit(int(params.Limit)).
		Offset(int((params.Page - 1) * params.Limit)).
		Order("created_at asc").
		Find(&comments).
		Error; err != nil {
		return nil, err
	}

	totalCount := len(comments)
	totalPages := (int64(len(comments)) + params.Limit - 1) / params.Limit

	return &dto.CommentListResponse{
		Comments: comments,
		Pagination: dto.PaginationResponse{
			Page:       params.Page,
			Limit:      params.Limit,
			TotalCount: int64(totalCount),
			TotalPages: totalPages,
		}}, nil
}

func (repo *CommentRepository) FindByID(ctx context.Context, id uint) (*models.Comment, error) {
	var comment models.Comment
	if err := repo.db.WithContext(ctx).First(&comment, id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (repo *CommentRepository) Update(ctx context.Context, comment *models.Comment) error {
	return repo.db.WithContext(ctx).Save(comment).Error
}

func (repo *CommentRepository) Delete(ctx context.Context, id uint) error {
	return repo.db.WithContext(ctx).Delete(&models.Comment{}, id).Error
}
