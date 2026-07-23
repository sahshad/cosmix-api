package services

import (
	"context"

	"post-service/internal/models"
	"post-service/internal/repositories"

	authEvents "cosmix/shared/events/auth"
	userEvents "cosmix/shared/events/user"

	"github.com/google/uuid"
)

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

func (svc *PostUserService) HandleUserUpdated(ctx context.Context, event userEvents.UserUpdated) error {
	postUser, err := svc.repo.FindByUserID(ctx, event.UserID)
	if err != nil {
		return err
	}

	postUser.Username = event.Username
	postUser.DisplayName = event.DisplayName
	postUser.AvatarURL = event.AvatarURL
	postUser.IsVerified = event.IsVerified
	postUser.UpdatedAt = &event.UpdatedAt

	return svc.repo.Update(ctx, postUser)
}

func (svc *PostUserService) FindByUserID(ctx context.Context, userID uuid.UUID) (*models.PostUser, error) {
	return svc.repo.FindByUserID(ctx, userID)
}

func (svc *PostUserService) Update(ctx context.Context, postUser *models.PostUser) error {
	return svc.repo.Update(ctx, postUser)
}

func (svc *PostUserService) Delete(ctx context.Context, userID uuid.UUID) error {
	return svc.repo.DeleteByID(ctx, userID)
}
