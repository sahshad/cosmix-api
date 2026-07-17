package repositories

import (
	"notification-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationPreferenceRepository struct {
	*BaseRepository[models.NotificationPreference]
}

func NewNotificationPreferenceRepository(
	db *gorm.DB,
) *NotificationPreferenceRepository {
	return &NotificationPreferenceRepository{
		NewBaseRepository[models.NotificationPreference](db),
	}
}

func (repo *NotificationPreferenceRepository) GetByUserID(userID uuid.UUID) (*models.NotificationPreference, error) {
	var preference models.NotificationPreference

	err := repo.db.
		Where("user_id = ?", userID).
		First(&preference).Error

	if err != nil {
		return nil, err
	}

	return &preference, nil
}
