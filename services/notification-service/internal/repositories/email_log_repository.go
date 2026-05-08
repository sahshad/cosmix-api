package repositories

import (
	"notification-service/internal/models"

	"gorm.io/gorm"
)

type EmailLogRepositoryInterface interface {
	Create(emailLog *models.EmailLog) error
	UpdateStatus(id uint, status string, errorMessage *string) error
}

type EmailLogRepository struct {
	db *gorm.DB
}

func NewEmailLogRepository(
	db *gorm.DB,
) EmailLogRepositoryInterface {
	return &EmailLogRepository{
		db: db,
	}
}

func (repo *EmailLogRepository) Create(emailLog *models.EmailLog) error {
	return repo.db.Create(emailLog).Error
}

func (repo *EmailLogRepository) UpdateStatus(id uint, status string, errorMessage *string) error {
	updates := map[string]any{
		"status": status,
	}

	if errorMessage != nil {
		updates["error_message"] = *errorMessage
	}

	return repo.db.
		Model(&models.EmailLog{}).
		Where("id = ?", id).
		Updates(updates).Error
}
