package services

import (
	"context"
	"errors"
	"post-service/internal/dto"
	"post-service/internal/models"
	"post-service/internal/repositories"
)

type PostServiceInterface interface {
	CreatePost(ctx context.Context, userID uint, req *dto.CreatePostRequest) (*models.Post, error)
	GetPostByID(ctx context.Context, id uint) (*models.Post, error)
	GetFeed(ctx context.Context, params *dto.PaginationRequest) (*dto.PostListResponse, error)
	GetUserPosts(ctx context.Context, userID uint, params *dto.PaginationRequest) (*dto.PostListResponse, error)
	UpdatePost(ctx context.Context, id uint, userID uint, req *dto.UpdatePostRequest) (*models.Post, error)
	DeletePost(ctx context.Context, id uint, userID uint) error
}

type PostService struct {
	repo repositories.PostRepositoryInterface
}

func NewPostService(repo repositories.PostRepositoryInterface) PostServiceInterface {
	return &PostService{repo: repo}
}

func (svc *PostService) CreatePost(ctx context.Context, userID uint, req *dto.CreatePostRequest) (*models.Post, error) {
	var media []models.PostMedia
	for _, m := range req.Media {
		media = append(media, models.PostMedia{
			PublicID: m.PublicID,
			URL:      m.URL,
			Type:     m.Type,
			Duration: &m.Duration,
		})
	}

	post := &models.Post{
		UserID:   userID,
		Content:  req.Content,
		Media:    media,
	}

	if err := svc.repo.Create(ctx, post); err != nil {
		return nil, err
	}
	return post, nil
}

func (svc *PostService) GetPostByID(ctx context.Context, id uint) (*models.Post, error) {
	return svc.repo.FindByID(ctx, id)
}

func (svc *PostService) GetFeed(ctx context.Context, params *dto.PaginationRequest) (*dto.PostListResponse, error) {
	return svc.repo.GetFeed(ctx, params)
}

func (svc *PostService) GetUserPosts(ctx context.Context, userID uint, params *dto.PaginationRequest) (*dto.PostListResponse, error) {
	return svc.repo.GetUserPosts(ctx, userID, params)
}

func (svc *PostService) UpdatePost(ctx context.Context, id uint, userID uint, req *dto.UpdatePostRequest) (*models.Post, error) {
	post, err := svc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if post.UserID != userID {
		return nil, errors.New("unauthorized to update this post")
	}

	post.Content = req.Content

	var media []models.PostMedia
	for _, m := range req.Media {
		media = append(media, models.PostMedia{
			PostID:   post.ID,
			PublicID: m.PublicID,
			URL:      m.URL,
			Type:     m.Type,
			Duration: &m.Duration,
		})
	}
	post.Media = media

	if err := svc.repo.Update(ctx, post); err != nil {
		return nil, err
	}
	return post, nil
}

func (svc *PostService) DeletePost(ctx context.Context, id uint, userID uint) error {
	post, err := svc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if post.UserID != userID {
		return errors.New("unauthorized to delete this post")
	}

	return svc.repo.Delete(ctx, id)
}
