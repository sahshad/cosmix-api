package repositories

import (
	"context"

	"post-service/internal/dto"
	"post-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentRepository struct {
	*BaseRepository[models.Comment]
}

func NewCommentRepository(
	db *gorm.DB,
) *CommentRepository {
	return &CommentRepository{
		NewBaseRepository[models.Comment](db),
	}
}

func (repo *CommentRepository) GetCommentsByPostID(ctx context.Context, postID uuid.UUID, viewerID uuid.UUID, params *dto.PaginationRequest) (*dto.CommentListResponse, error) {
	var comments []models.Comment

	var totalCount int64
	err := repo.db.WithContext(ctx).
		Model(&models.Comment{}).
		Where("post_id = ? AND parent_comment_id IS NULL", postID).
		Count(&totalCount).Error
	if err != nil {
		return nil, err
	}

	if err := repo.db.WithContext(ctx).
		Model(&models.Comment{}).
		Preload("User").
		Where("post_id = ? AND parent_comment_id IS NULL", postID).
		Limit(int(params.Limit)).
		Offset(int((params.Page - 1) * params.Limit)).
		Order("created_at asc").
		Find(&comments).
		Error; err != nil {
		return nil, err
	}

	return repo.buildCommentListResponse(ctx, comments, viewerID, params, int32(totalCount))
}

func (repo *CommentRepository) GetReplies(ctx context.Context, parentCommentID uuid.UUID, viewerID uuid.UUID, params *dto.PaginationRequest) (*dto.CommentListResponse, error) {
	var comments []models.Comment

	var totalCount int64
	err := repo.db.WithContext(ctx).
		Model(&models.Comment{}).
		Where("parent_comment_id = ?", parentCommentID).
		Count(&totalCount).Error
	if err != nil {
		return nil, err
	}

	if err := repo.db.WithContext(ctx).
		Model(&models.Comment{}).
		Preload("User").
		Where("parent_comment_id = ?", parentCommentID).
		Limit(int(params.Limit)).
		Offset(int((params.Page - 1) * params.Limit)).
		Order("created_at asc").
		Find(&comments).
		Error; err != nil {
		return nil, err
	}

	return repo.buildCommentListResponse(ctx, comments, viewerID, params, int32(totalCount))
}

func (repo *CommentRepository) buildCommentListResponse(ctx context.Context, comments []models.Comment, viewerID uuid.UUID, params *dto.PaginationRequest, totalCount int32) (*dto.CommentListResponse, error) {
	likedCommentIDs, err := repo.getLikedCommentIDs(ctx, viewerID, comments)
	if err != nil {
		return nil, err
	}

	totalPages := (totalCount + params.Limit - 1) / params.Limit

	commentList := make([]dto.Comment, 0, len(comments))
	for _, comment := range comments {
		commentList = append(commentList, toComment(comment, likedCommentIDs[comment.ID], comment.UserID == viewerID))
	}

	return &dto.CommentListResponse{
		Comments: commentList,
		Pagination: dto.PaginationResponse{
			Page:       params.Page,
			Limit:      params.Limit,
			TotalCount: totalCount,
			TotalPages: totalPages,
		}}, nil
}

// getLikedCommentIDs returns the subset of comment IDs the viewer has liked.
// When viewerID is uuid.Nil (anonymous request), it returns an empty set.
func (repo *CommentRepository) getLikedCommentIDs(ctx context.Context, viewerID uuid.UUID, comments []models.Comment) (map[uuid.UUID]bool, error) {
	liked := make(map[uuid.UUID]bool)

	if viewerID == uuid.Nil || len(comments) == 0 {
		return liked, nil
	}

	commentIDs := make([]uuid.UUID, len(comments))
	for i, comment := range comments {
		commentIDs[i] = comment.ID
	}

	var likedCommentIDs []uuid.UUID
	if err := repo.db.WithContext(ctx).
		Model(&models.CommentLike{}).
		Where("user_id = ? AND comment_id IN ?", viewerID, commentIDs).
		Pluck("comment_id", &likedCommentIDs).
		Error; err != nil {
		return nil, err
	}

	for _, id := range likedCommentIDs {
		liked[id] = true
	}

	return liked, nil
}

func toComment(comment models.Comment, isLiked bool, isOwner bool) dto.Comment {
	avatarURL := ""
	if comment.User.AvatarURL != nil {
		avatarURL = *comment.User.AvatarURL
	}

	return dto.Comment{
		ID:              comment.ID,
		PostID:          comment.PostID,
		AuthorID:        comment.UserID,
		ParentCommentID: comment.ParentCommentID,
		Content:         comment.Content,
		LikesCount:      comment.LikesCount,
		RepliesCount:    comment.RepliesCount,
		IsLiked:         isLiked,
		IsOwner:         isOwner,
		CreatedAt:       formatTime(comment.CreatedAt),
		UpdatedAt:       formatTimePtr(comment.UpdatedAt),
		User: dto.User{
			ID:          comment.User.UserID,
			Username:    comment.User.Username,
			DisplayName: comment.User.DisplayName,
			AvatarURL:   avatarURL,
			CreatedAt:   formatTime(comment.User.CreatedAt),
			UpdatedAt:   formatTimePtr(comment.User.UpdatedAt),
		},
	}
}

func (repo *CommentRepository) FindCommentByID(ctx context.Context, id uuid.UUID, viewerID uuid.UUID) (*dto.Comment, error) {
	var comment models.Comment

	if err := repo.db.WithContext(ctx).
		Model(&models.Comment{}).
		Preload("User").
		Where("id = ?", id).
		First(&comment).Error; err != nil {
		return nil, err
	}

	likedCommentIDs, err := repo.getLikedCommentIDs(ctx, viewerID, []models.Comment{comment})
	if err != nil {
		return nil, err
	}

	result := toComment(comment, likedCommentIDs[comment.ID], comment.UserID == viewerID)
	return &result, nil
}

func (repo *CommentRepository) UpdateContent(ctx context.Context, id uuid.UUID, content string) error {
	return repo.db.WithContext(ctx).
		Model(&models.Comment{}).
		Where("id = ?", id).
		Update("content", content).Error
}

func (repo *CommentRepository) IncrementLikesCount(ctx context.Context, id uuid.UUID, delta int) error {
	return repo.db.WithContext(ctx).
		Model(&models.Comment{}).
		Where("id = ?", id).
		UpdateColumn("likes_count", gorm.Expr("likes_count + ?", delta)).
		Error
}

func (repo *CommentRepository) IncrementRepliesCount(ctx context.Context, id uuid.UUID, delta int) error {
	return repo.db.WithContext(ctx).
		Model(&models.Comment{}).
		Where("id = ?", id).
		UpdateColumn("replies_count", gorm.Expr("replies_count + ?", delta)).
		Error
}

// CountDescendants returns the number of comments in the reply subtree
// rooted under commentID (i.e. everything that will be cascade-deleted
// along with it), not including commentID itself.
func (repo *CommentRepository) CountDescendants(ctx context.Context, commentID uuid.UUID) (int, error) {
	var count int64
	err := repo.db.WithContext(ctx).Raw(`
		WITH RECURSIVE descendants AS (
			SELECT id FROM comments WHERE parent_comment_id = ?
			UNION ALL
			SELECT c.id FROM comments c
			INNER JOIN descendants d ON c.parent_comment_id = d.id
		)
		SELECT count(*) FROM descendants
	`, commentID).Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (repo *CommentRepository) GetPostByID(ctx context.Context, postID uuid.UUID) (*models.Post, error) {
	var post models.Post
	if err := repo.db.WithContext(ctx).
		Model(&models.Post{}).
		Where("id = ?", postID).
		First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}
