package services

import (
	"context"
	authEvents "cosmix/shared/events/auth"
	"errors"
	"time"
	"user-service/internal/dto"
	"user-service/internal/models"
	"user-service/internal/repositories"
)

type UserProfileServiceInterface interface {
	GetProfile(ctx context.Context, userID uint) (*dto.UserProfileResponse, error)
	GetProfileByID(ctx context.Context, id uint) (*dto.UserProfileResponse, error)
	GetProfileByUsername(ctx context.Context, username string) (*dto.UserProfileResponse, error)
	UpdateProfile(ctx context.Context, userID uint, input dto.UpdateProfileDTO) (*dto.UserProfileResponse, error)
	CreateProfile(ctx context.Context, profile *models.UserProfile) error
	CreateFromAuthEvent(event authEvents.AuthUserRegistered) error
}

type UserProfileService struct {
	repo repositories.UserProfileRepositoryInterface
}

func NewUserProfileService(repo repositories.UserProfileRepositoryInterface) UserProfileServiceInterface {
	return &UserProfileService{
		repo: repo,
	}
}

func (svc *UserProfileService) GetProfile(ctx context.Context, userID uint) (*dto.UserProfileResponse, error) {
	profile, err := svc.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("profile not found")
	}
	return svc.toResponse(profile), nil
}

func (svc *UserProfileService) GetProfileByID(ctx context.Context, id uint) (*dto.UserProfileResponse, error) {
	profile, err := svc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, errors.New("profile not found")
	}
	return svc.toResponse(profile), nil
}

func (svc *UserProfileService) UpdateProfile(ctx context.Context, userID uint, input dto.UpdateProfileDTO) (*dto.UserProfileResponse, error) {
	profile, err := svc.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("profile not found")
	}

	if input.FirstName != nil {
		profile.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		profile.LastName = *input.LastName
	}
	if input.Username != nil {
		profile.Username = input.Username
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
	if input.Bio != nil {
		profile.Bio = input.Bio
	}

	if err := svc.repo.Update(ctx, profile); err != nil {
		return nil, err
	}

	return svc.toResponse(profile), nil
}

func (svc *UserProfileService) CreateProfile(ctx context.Context, profile *models.UserProfile) error {
	return svc.repo.Create(ctx, profile)
}

func (svc *UserProfileService) CreateFromAuthEvent(event authEvents.AuthUserRegistered) error {
	ctx := context.Background()
	profile := &models.UserProfile{
		AuthUserID: event.AuthUserID,
		Email:      event.Email,
		FirstName:  event.FirstName,
		LastName:   event.LastName,
		CreatedAt:  event.CreatedAt,
	}
	return svc.repo.Create(ctx, profile)
}

func (svc *UserProfileService) GetProfileByUsername(ctx context.Context, username string) (*dto.UserProfileResponse, error) {
	profile, err := svc.repo.FindByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("profile not found")
	}
	return svc.toResponse(profile), nil
}

func (svc *UserProfileService) toResponse(profile *models.UserProfile) *dto.UserProfileResponse {
	return &dto.UserProfileResponse{
		User: dto.UserResponse{
			ID:          profile.ID,
			AuthUserID:  profile.AuthUserID,
			FirstName:   profile.FirstName,
			LastName:    profile.LastName,
			Username:    profile.Username,
			DateOfBirth: profile.DateOfBirth,
			AvatarURL:   profile.AvatarURL,
			Bio:         profile.Bio,
			CreatedAt:   profile.CreatedAt,
			UpdatedAt:   profile.UpdatedAt,
		},
	}
}
