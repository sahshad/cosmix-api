package services

import (
	"context"
	"errors"
	"time"

	"user-service/internal/dto"
	"user-service/internal/models"
	"user-service/internal/repositories"

	"cosmix/shared/core/eventbus"
	"cosmix/shared/core/rabbitmq"

	authEvents "cosmix/shared/events/auth"
	userEvents "cosmix/shared/events/user"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type UserService struct {
	repo     *repositories.UserRepository
	rabbitCh *amqp.Channel
}

func NewUserService(
	repo *repositories.UserRepository,
	rabbitCh *amqp.Channel,
) *UserService {
	return &UserService{
		repo:     repo,
		rabbitCh: rabbitCh,
	}
}

func (svc *UserService) GetProfile(ctx context.Context, userID uuid.UUID) (*dto.UserProfileResponse, error) {
	profile, err := svc.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("profile not found")
	}
	return svc.toResponse(profile), nil
}

func (svc *UserService) GetProfileByID(ctx context.Context, id uuid.UUID) (*dto.UserProfileResponse, error) {
	profile, err := svc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("profile not found")
	}
	return svc.toResponse(profile), nil
}

func (svc *UserService) UpdateProfile(ctx context.Context, userID uuid.UUID, input dto.UpdateProfileDTO) (*dto.UserProfileResponse, error) {
	profile, err := svc.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("profile not found")
	}

	if input.DisplayName != nil {
		profile.DisplayName = *input.DisplayName
	}

	if input.Username != nil {
		profile.Username = *input.Username
	}

	if input.DateOfBirth != nil {
		dob, err := time.Parse("2006-01-02", *input.DateOfBirth)
		if err != nil {
			return nil, errors.New("invalid date of birth format")
		}
		profile.DateOfBirth = &dob
	}

	if input.AvatarURL != nil {
		profile.AvatarURL = input.AvatarURL
	}

	if input.CoverImageURL != nil {
		profile.CoverImageURL = input.CoverImageURL
	}

	if input.Website != nil {
		profile.Website = input.Website
	}

	if input.Location != nil {
		profile.Location = input.Location
	}

	if input.Bio != nil {
		profile.Bio = input.Bio
	}

	now := time.Now().UTC()
	profile.UpdatedAt = &now

	if err := svc.repo.Update(ctx, profile); err != nil {
		return nil, err
	}

	requestID := eventbus.RequestID(ctx)
	ctx = eventbus.WithCorrelationID(ctx, requestID)

	event := userEvents.UserUpdated{
		EventVersion: userEvents.EventVersionOne,
		UserID:       profile.UserID,
		Username:     profile.Username,
		DisplayName:  profile.DisplayName,
		AvatarURL:    profile.AvatarURL,
		IsVerified:   profile.IsVerified,
		UpdatedAt:    now,
	}

	// @Publish - user profile updated
	if err := eventbus.Publish(
		ctx,
		svc.rabbitCh,
		rabbitmq.ExchangeEvents,
		rabbitmq.UserProfileUpdated,
		event,
	); err != nil {
		return nil, err
	}

	return svc.toResponse(profile), nil
}

func (svc *UserService) CreateProfile(ctx context.Context, profile *models.User) error {
	return svc.repo.Create(ctx, profile)
}

func (svc *UserService) HandleAuthUserEmailVerificationCompleted(ctx context.Context, event authEvents.AuthUserEmailVerificationCompleted) error {

	profile := &models.User{
		UserID:      event.AuthUserID,
		Email:       event.Email,
		DisplayName: event.DisplayName,
		Username:    event.Username,
		CreatedAt:   event.CreatedAt,
	}
	return svc.repo.Create(ctx, profile)
}

func (svc *UserService) GetProfileByUsername(ctx context.Context, username string) (*dto.UserProfileResponse, error) {
	profile, err := svc.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("profile not found")
	}
	return svc.toResponse(profile), nil
}

func (svc *UserService) toResponse(profile *models.User) *dto.UserProfileResponse {
	return &dto.UserProfileResponse{
		User: dto.UserResponse{
			UserID:        profile.UserID,
			DisplayName:   profile.DisplayName,
			Username:      profile.Username,
			Email:         profile.Email,
			IsPrivate:     profile.IsPrivate,
			IsVerified:    profile.IsVerified,
			IsActive:      profile.IsActive,
			DateOfBirth:   profile.DateOfBirth,
			AvatarURL:     profile.AvatarURL,
			CoverImageURL: profile.CoverImageURL,
			Website:       profile.Website,
			Location:      profile.Location,
			Bio:           profile.Bio,
			LastSeenAt:    profile.LastSeenAt,
			CreatedAt:     profile.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     formatTimePtr(profile.UpdatedAt),
		},
	}
}

func formatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}
