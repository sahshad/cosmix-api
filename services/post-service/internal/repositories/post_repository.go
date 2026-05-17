package repositories

import (
	"context"
	"post-service/internal/dto"
	"post-service/internal/models"

	"gorm.io/gorm"
)

type PostRepositoryInterface interface {
	Create(ctx context.Context, post *models.Post) error
	FindByID(ctx context.Context, id uint) (*models.Post, error)
	GetFeed(ctx context.Context, params *dto.PaginationRequest) (*dto.PostListResponse, error)
	GetUserPosts(ctx context.Context, userID uint, params *dto.PaginationRequest) (*dto.PostListResponse, error)
	Update(ctx context.Context, post *models.Post) error
	Delete(ctx context.Context, id uint) error
}

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepositoryInterface {
	return &PostRepository{db: db}
}

func (repo *PostRepository) Create(ctx context.Context, post *models.Post) error {
	return repo.db.WithContext(ctx).Create(post).Error
}

func (repo *PostRepository) FindByID(ctx context.Context, id uint) (*models.Post, error) {
	var post models.Post
	if err := repo.db.WithContext(ctx).Preload("Media").Preload("Likes").Preload("Comments").First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (repo *PostRepository) GetFeed(ctx context.Context, params *dto.PaginationRequest) (*dto.PostListResponse, error) {
	var posts []models.Post
	if err := repo.db.WithContext(ctx).
		Table("posts").
		Preload("User").
		Preload("Media").
		Order("created_at desc").
		Limit(int(params.Limit)).
		Offset(int((params.Page - 1) * params.Limit)).
		Find(&posts).
		Error; err != nil {
		return nil, err
	}

	totalCount := int64(len(posts))
	totalPages := (totalCount + params.Limit - 1) / params.Limit

	postList := make([]dto.PostList, len(posts))
	for i, post := range posts {
		var mediaList []dto.Media
		for _, media := range post.Media {
			mediaList = append(mediaList, dto.Media{
				ID:        media.ID,
				Type:      media.Type,
				URL:       media.URL,
				PublicID:  media.PublicID,
				Duration:  media.Duration,
				CreatedAt: media.CreatedAt,
			})
		}
		postList[i] = dto.PostList{
			ID:        post.ID,
			Content:   post.Content,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
			User: dto.User{
				ID:          post.User.UserID,
				Username:    post.User.Username,
				DisplayName: post.User.DisplayName,
				CreatedAt:   post.User.CreatedAt,
				UpdatedAt:   post.User.UpdatedAt,
			},
			Media: mediaList,
		}
	}
	return &dto.PostListResponse{
		Posts: postList,
		Pagination: dto.PaginationResponse{
			Page:       params.Page,
			Limit:      params.Limit,
			TotalCount: totalCount,
			TotalPages: totalPages,
		}}, nil
}

func (repo *PostRepository) GetUserPosts(ctx context.Context, userID uint, params *dto.PaginationRequest) (*dto.PostListResponse, error) {
	var posts []dto.PostList
	if err := repo.db.WithContext(ctx).
		Table("posts").
		Preload("User").
		Preload("Media").
		// Preload("Likes").
		// Preload("Comments").
		Where("user.user_id = ?", userID).
		Order("posts.created_at desc").
		Limit(int(params.Limit)).
		Offset(int((params.Page - 1) * params.Limit)).
		Find(&posts).
		Error; err != nil {
		return nil, err
	}

	totalCount := int64(len(posts))
	totalPages := (totalCount + params.Limit - 1) / params.Limit

	return &dto.PostListResponse{
		Posts: posts,
		Pagination: dto.PaginationResponse{
			Page:       params.Page,
			Limit:      params.Limit,
			TotalCount: totalCount,
			TotalPages: totalPages,
		}}, nil
}

func (repo *PostRepository) Update(ctx context.Context, post *models.Post) error {
	if err := repo.db.WithContext(ctx).Where("post_id = ?", post.ID).Delete(&models.PostMedia{}).Error; err != nil {
		return err
	}
	return repo.db.WithContext(ctx).Save(post).Error
}

func (repo *PostRepository) Delete(ctx context.Context, id uint) error {
	return repo.db.WithContext(ctx).Delete(&models.Post{}, id).Error
}
