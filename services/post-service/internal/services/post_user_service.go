package services

import (
	"context"

	"post-service/internal/models"
	"post-service/internal/repositories"

	authEvents "cosmix/shared/events/auth"
)

// type PostUserServiceInterface interface {
// 	CreateFromAuthEvent(ctx context.Context, event authEvents.AuthUserRegistered) error
// 	FindByUserID(ctx context.Context, userID uint) (*models.PostUser, error)
// 	Update(ctx context.Context, postUser *models.PostUser) error
// 	Delete(ctx context.Context, userID uint) error
// }

type PostUserService struct {
	repo *repositories.PostUserRepository
}

func NewPostUserService(
	repo *repositories.PostUserRepository,
) *PostUserService {
	return &PostUserService{
		repo: repo,
	}
}

func (svc *PostUserService) HandleAuthUserEmailVerificationCompleted(ctx context.Context, event authEvents.AuthUserEmailVerificationCompleted) error {
	postUser := &models.PostUser{
		UserID:      event.AuthUserID,
		Username:    event.Username,
		DisplayName: event.DisplayName,
		CreatedAt:   event.CreatedAt,
	}

	if _, err := svc.repo.Create(ctx, postUser); err != nil {
		return err
	}

	return nil
}

func (svc *PostUserService) FindByUserID(ctx context.Context, userID uint) (*models.PostUser, error) {
	return svc.repo.FindByUserID(ctx, userID)
}

func (svc *PostUserService) Update(ctx context.Context, postUser *models.PostUser) error {
	return svc.repo.Update(ctx, postUser)
}

func (svc *PostUserService) Delete(ctx context.Context, userID uint) error {
	return svc.repo.DeleteByID(ctx, userID)
}
