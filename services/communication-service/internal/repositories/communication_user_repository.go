package repositories

import (
	"context"

	"communication-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CommunicationUserRepository struct {
	*BaseRepository[models.CommunicationUser]
}

func NewCommunicationUserRepository(
	db *gorm.DB,
) *CommunicationUserRepository {
	return &CommunicationUserRepository{
		NewBaseRepository[models.CommunicationUser](db),
	}
}

func (repo *CommunicationUserRepository) Create(ctx context.Context, user *models.CommunicationUser) error {
	return repo.db.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			DoNothing: true,
		}).
		Create(user).
		Error
}

func (repo *CommunicationUserRepository) Upsert(ctx context.Context, user *models.CommunicationUser) error {
	return repo.db.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"username",
				"display_name",
				"avatar_url",
				"is_verified",
				"updated_at",
			}),
		}).
		Create(user).
		Error
}

func (repo *CommunicationUserRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*models.CommunicationUser, error) {
	var user models.CommunicationUser
	if err := repo.db.
		WithContext(ctx).
		Where("user_id = ?", userID).
		First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
