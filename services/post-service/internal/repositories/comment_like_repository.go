package repositories

import (
	"context"

	"post-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentLikeRepository struct {
	*BaseRepository[models.CommentLike]
}

func NewCommentLikeRepository(
	db *gorm.DB,
) *CommentLikeRepository {
	return &CommentLikeRepository{
		NewBaseRepository[models.CommentLike](db),
	}
}

func (repo *CommentLikeRepository) Delete(ctx context.Context, commentID uuid.UUID, userID uuid.UUID) (int64, error) {
	tx := repo.db.WithContext(ctx).Where("comment_id = ? AND user_id = ?", commentID, userID).Delete(&models.CommentLike{})
	return tx.RowsAffected, tx.Error
}

func (repo *CommentLikeRepository) Exists(ctx context.Context, commentID uuid.UUID, userID uuid.UUID) (bool, error) {
	var count int64
	err := repo.db.WithContext(ctx).Model(&models.CommentLike{}).Where("comment_id = ? AND user_id = ?", commentID, userID).Count(&count).Error
	return count > 0, err
}
