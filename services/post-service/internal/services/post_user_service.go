package services

import (
	"context"
	"post-service/internal/models"
	"post-service/internal/repositories"

	authEvents "cosmix/shared/events/auth"
)

type PostUserServiceInterface interface {
	CreateFromAuthEvent(event authEvents.AuthUserRegistered) error
	FindByUserID(ctx context.Context, userID uint) (*models.PostUser, error)
	Update(ctx context.Context, postUser *models.PostUser) error
	Delete(ctx context.Context, userID uint) error
}

type PostUserService struct {
	postUserRepository repositories.PostUserRepositoryInterface
}

func NewPostUserService(
	postUserRepository repositories.PostUserRepositoryInterface,
) PostUserServiceInterface {
	return &PostUserService{
		postUserRepository: postUserRepository,
	}
}

func (svc *PostUserService) CreateFromAuthEvent(event authEvents.AuthUserRegistered) error {
	ctx := context.Background()
	postUser := &models.PostUser{
		UserID:      event.AuthUserID,
		Username:    event.Username,
		DisplayName: event.DisplayName,
		CreatedAt:   event.CreatedAt,
	}

	if _, err := svc.postUserRepository.Create(ctx, postUser); err != nil {
		return err
	}

	return nil
}

func (svc *PostUserService) FindByUserID(ctx context.Context, userID uint) (*models.PostUser, error) {
	return svc.postUserRepository.FindByUserID(ctx, userID)
}

func (svc *PostUserService) Update(ctx context.Context, postUser *models.PostUser) error {
	return svc.postUserRepository.Update(ctx, postUser)
}

func (svc *PostUserService) Delete(ctx context.Context, userID uint) error {
	return svc.postUserRepository.Delete(ctx, userID)
}
