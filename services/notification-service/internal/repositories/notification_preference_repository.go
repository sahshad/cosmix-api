package repositories

import (
	"notification-service/internal/models"

	"gorm.io/gorm"
)

// type NotificationPreferenceRepositoryInterface interface {
// 	Create(preference *models.NotificationPreference) error
// 	GetByUserID(userID uint) (*models.NotificationPreference, error)
// 	Update(preference *models.NotificationPreference) error
// }

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

// func (repo *NotificationPreferenceRepository) Create(preference *models.NotificationPreference) error {
// 	return repo.db.Create(preference).Error
// }

func (repo *NotificationPreferenceRepository) GetByUserID(userID uint) (*models.NotificationPreference, error) {
	var preference models.NotificationPreference

	err := repo.db.
		Where("user_id = ?", userID).
		First(&preference).Error

	if err != nil {
		return nil, err
	}

	return &preference, nil
}

// func (repo *NotificationPreferenceRepository) Update(preference *models.NotificationPreference) error {
// 	return repo.db.Save(preference).Error
// }
