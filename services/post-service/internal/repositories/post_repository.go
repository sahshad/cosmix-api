package repositories

import (
	"context"
	"post-service/internal/dto"
	"post-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// type PostRepositoryInterface interface {
// 	Create(ctx context.Context, post *models.Post) error
// 	FindByID(ctx context.Context, id uuid.UUID) (*models.Post, error)
// 	GetFeed(ctx context.Context, params *dto.PaginationRequest) (*dto.PostListResponse, error)
// 	GetUserPosts(ctx context.Context, userID uuid.UUID, params *dto.PaginationRequest) (*dto.PostListResponse, error)
// 	Update(ctx context.Context, post *models.Post) error
// 	Delete(ctx context.Context, id uuid.UUID) error
// }

type PostRepository struct {
	*BaseRepository[models.Post]
}

func NewPostRepository(
	db *gorm.DB,
) *PostRepository {
	return &PostRepository{
		NewBaseRepository[models.Post](db),
	}
}

// func (repo *PostRepository) Create(ctx context.Context, post *models.Post) error {
// 	return repo.db.WithContext(ctx).Create(post).Error
// }

func (repo *PostRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Post, error) {
	var post models.Post
	if err := repo.db.WithContext(ctx).Preload("Media").Preload("Likes").Preload("Comments").First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (repo *PostRepository) GetFeed(ctx context.Context, viewerID uuid.UUID, params *dto.PaginationRequest) (*dto.PostListResponse, error) {
	var posts []models.Post
	if err := repo.db.WithContext(ctx).
		Model(&models.Post{}).
		Preload("User").
		Preload("Media").
		Order("created_at desc").
		Limit(int(params.Limit)).
		Offset(int((params.Page - 1) * params.Limit)).
		Find(&posts).
		Error; err != nil {
		return nil, err
	}

	likedPostIDs, err := repo.getLikedPostIDs(ctx, viewerID, posts)
	if err != nil {
		return nil, err
	}

	totalCount := int32(len(posts))
	totalPages := (totalCount + params.Limit - 1) / params.Limit

	postList := make([]dto.PostList, 0, len(posts))
	for _, post := range posts {
		postList = append(postList, toPostList(post, likedPostIDs[post.ID], post.UserID == viewerID))
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

// getLikedPostIDs returns the subset of postIDs the viewer has liked. When
// viewerID is uuid.Nil (anonymous request), it returns an empty set.
func (repo *PostRepository) getLikedPostIDs(ctx context.Context, viewerID uuid.UUID, posts []models.Post) (map[uuid.UUID]bool, error) {
	liked := make(map[uuid.UUID]bool)

	if viewerID == uuid.Nil || len(posts) == 0 {
		return liked, nil
	}

	postIDs := make([]uuid.UUID, len(posts))
	for i, post := range posts {
		postIDs[i] = post.ID
	}

	var likedPostIDs []uuid.UUID
	if err := repo.db.WithContext(ctx).
		Model(&models.Like{}).
		Where("user_id = ? AND post_id IN ?", viewerID, postIDs).
		Pluck("post_id", &likedPostIDs).
		Error; err != nil {
		return nil, err
	}

	for _, id := range likedPostIDs {
		liked[id] = true
	}

	return liked, nil
}

func (repo *PostRepository) GetUserPosts(ctx context.Context, userID uuid.UUID, viewerID uuid.UUID, params *dto.PaginationRequest) (*dto.PostListResponse, error) {
	var posts []models.Post
	if err := repo.db.WithContext(ctx).
		Model(&models.Post{}).
		Preload("User").
		Preload("Media").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Limit(int(params.Limit)).
		Offset(int((params.Page - 1) * params.Limit)).
		Find(&posts).
		Error; err != nil {
		return nil, err
	}

	likedPostIDs, err := repo.getLikedPostIDs(ctx, viewerID, posts)
	if err != nil {
		return nil, err
	}

	totalCount := int32(len(posts))
	totalPages := (totalCount + params.Limit - 1) / params.Limit

	postList := make([]dto.PostList, 0, len(posts))
	for _, post := range posts {
		postList = append(postList, toPostList(post, likedPostIDs[post.ID], post.UserID == viewerID))
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

func toPostList(post models.Post, isLiked bool, isOwner bool) dto.PostList {
	mediaList := make([]dto.Media, 0, len(post.Media))

	avatarURL := ""
	if post.User.AvatarURL != nil {
		avatarURL = *post.User.AvatarURL
	}

	for _, media := range post.Media {
		mediaList = append(mediaList, dto.Media{
			ID:        media.ID,
			PostID:    media.PostID,
			PublicID:  media.PublicID,
			URL:       media.URL,
			Type:      media.Type,
			Duration:  media.Duration,
			CreatedAt: formatTime(media.CreatedAt),
		})
	}

	return dto.PostList{
		ID:            post.ID,
		Content:       post.Content,
		LikesCount:    post.LikesCount,
		CommentsCount: post.CommentsCount,
		IsLiked:       isLiked,
		IsOwner:       isOwner,
		CreatedAt:     formatTime(post.CreatedAt),
		UpdatedAt:     formatTimePtr(post.UpdatedAt),

		User: dto.User{
			ID:          post.User.UserID,
			Username:    post.User.Username,
			DisplayName: post.User.DisplayName,
			AvatarURL:   avatarURL,
			CreatedAt:   formatTime(post.User.CreatedAt),
			UpdatedAt:   formatTimePtr(post.User.UpdatedAt),
		},

		Media: mediaList,
	}
}

func (repo *PostRepository) Update(ctx context.Context, post *models.Post) error {
	if err := repo.db.WithContext(ctx).Where("post_id = ?", post.ID).Delete(&models.PostMedia{}).Error; err != nil {
		return err
	}
	return repo.db.WithContext(ctx).Save(post).Error
}

func (repo *PostRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return repo.db.WithContext(ctx).Delete(&models.Post{}, id).Error
}

func (repo *PostRepository) IncrementLikesCount(ctx context.Context, id uuid.UUID, delta int) error {
	return repo.db.WithContext(ctx).
		Model(&models.Post{}).
		Where("id = ?", id).
		UpdateColumn("likes_count", gorm.Expr("likes_count + ?", delta)).
		Error
}

func (repo *PostRepository) IncrementCommentsCount(ctx context.Context, id uuid.UUID, delta int) error {
	return repo.db.WithContext(ctx).
		Model(&models.Post{}).
		Where("id = ?", id).
		UpdateColumn("comments_count", gorm.Expr("comments_count + ?", delta)).
		Error
}
