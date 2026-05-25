package repositories

import (
	"auth-service/internal/models"
	"context"

	"gorm.io/gorm"
)

// type UserRepositoryInterface interface {
// 	Create(ctx context.Context, u *models.AuthUser) error
// 	FindByEmail(ctx context.Context, email string) (*models.AuthUser, error)
// 	FindByID(ctx context.Context, id uint) (*models.AuthUser, error)
// 	FindByUsername(ctx context.Context, username string) (*models.AuthUser, error)
// 	Update(ctx context.Context, u *models.AuthUser) error
// }

type AuthUserRepository struct {
	*BaseRepository[models.AuthUser]
}

func NewAuthUserRepository(
	db *gorm.DB,
) *AuthUserRepository {
	return &AuthUserRepository{
		NewBaseRepository[models.AuthUser](db),
	}
}

// func (repo *UserRepo) Create(ctx context.Context, u *models.AuthUser) error {
// 	return repo.db.WithContext(ctx).Create(u).Error
// }

func (repo *AuthUserRepository) FindByEmail(ctx context.Context, email string) (*models.AuthUser, error) {
	var user models.AuthUser
	if err := repo.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// func (repo *UserRepo) FindByID(ctx context.Context, id uint) (*models.AuthUser, error) {
// 	var user models.AuthUser
// 	if err := repo.db.WithContext(ctx).First(&user, id).Error; err != nil {
// 		return nil, err
// 	}
// 	return &user, nil
// }

func (repo *AuthUserRepository) FindByUsername(ctx context.Context, username string) (*models.AuthUser, error) {
	var user models.AuthUser
	if err := repo.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// func (repo *UserRepo) Update(ctx context.Context, u *models.AuthUser) error {
// 	return repo.db.WithContext(ctx).Save(u).Error
// }
