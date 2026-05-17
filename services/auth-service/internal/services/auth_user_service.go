package services

import (
	"context"
	"errors"
	"time"

	"auth-service/internal/dto"
	"auth-service/internal/models"
	"auth-service/internal/repositories"
	"auth-service/internal/utils"
	appErr "cosmix/shared/core/errors"

	"gorm.io/gorm"
)

type AuthServiceInterface interface {
	Register(ctx context.Context, input dto.RegisterDTO) (*models.AuthUser, error)
	Login(ctx context.Context, input dto.LoginDTO) (*dto.LoginResponseDTO, error)
	Refresh(ctx context.Context, refreshToken string) (*dto.RefreshResponseDTO, error)
	GetByID(ctx context.Context, id uint) (*models.AuthUser, error)
	UpdateUserPassword(ctx context.Context, userID uint, newPassword string) error
}

type AuthService struct {
	userRepo    repositories.UserRepositoryInterface
	sessionRepo repositories.AuthSessionRepositoryInterface
}

func NewAuthService(
	userRepo repositories.UserRepositoryInterface,
	sessionRepo repositories.AuthSessionRepositoryInterface,
) AuthServiceInterface {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (svc *AuthService) Register(ctx context.Context, input dto.RegisterDTO) (*models.AuthUser, error) {
	if _, err := svc.userRepo.FindByEmail(ctx, input.Email); err == nil {
		return nil, appErr.NewBadRequest("email", "Email already in use")
	}

	pwHash, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := &models.AuthUser{
		Email:        input.Email,
		PasswordHash: string(pwHash),
	}

	if err := svc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (svc *AuthService) Login(ctx context.Context, input dto.LoginDTO) (*dto.LoginResponseDTO, error) {
	authUser, err := svc.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, appErr.NewBadRequest("email", "Invalid credentials")
	}

	if err := utils.VerifyPassword(input.Password, authUser.PasswordHash); err != nil {
		return nil, appErr.NewBadRequest("password", "Invalid credentials")
	}

	access, err := utils.GenerateAccessToken(authUser.ID, "user")
	if err != nil {
		return nil, err
	}
	refresh, err := utils.GenerateRefreshToken(authUser.ID, "user")
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	authUser.LastLoginAt = &now
	if err := svc.userRepo.Update(ctx, authUser); err != nil {
		return nil, err
	}

	tokenHash, err := utils.HashPassword(refresh)
	if err != nil {
		return nil, err
	}

	session := models.AuthSession{
		AuthUserID:       authUser.ID,
		RefreshTokenHash: string(tokenHash),
		Device:           "",
		ExpiresAt:        time.Now().AddDate(0, 1, 0),
		CreatedAt:        time.Now().UTC(),
	}

	if _, err := svc.sessionRepo.Create(&session); err != nil {
		return nil, err
	}

	return &dto.LoginResponseDTO{
		AccessToken:  access,
		RefreshToken: refresh,
		AuthUser: &dto.AuthUserDTO{
			Email:         authUser.Email,
			IsActive:      authUser.IsActive,
			EmailVerified: authUser.EmailVerified,
			LastLoginAt:   authUser.LastLoginAt,
			CreatedAt:     authUser.CreatedAt,
			UpdatedAt:     authUser.UpdatedAt,
		},
	}, nil
}

func (svc *AuthService) Refresh(ctx context.Context, refreshToken string) (*dto.RefreshResponseDTO, error) {
	claims, err := utils.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, appErr.NewUnauthorized("Invalid Token")
	}

	if _, err := svc.userRepo.FindByID(ctx, claims.UserID); err != nil {
		return nil, appErr.NewUnauthorized("Invalid Token")
	}

	newAccess, err := utils.GenerateAccessToken(claims.UserID, "user")
	if err != nil {
		return nil, err
	}
	newRefresh, err := utils.GenerateRefreshToken(claims.UserID, "user")
	if err != nil {
		return nil, err
	}
	return &dto.RefreshResponseDTO{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
	}, nil
}

func (svc *AuthService) GetByID(ctx context.Context, id uint) (*models.AuthUser, error) {
	user, err := svc.userRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErr.NewNotFound("user")
		}
		return nil, err
	}
	return user, nil
}

func (svc *AuthService) UpdateUserPassword(ctx context.Context, userID uint, newPassword string) error {
	user, err := svc.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	pwHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	now := time.Now().UTC()

	user.PasswordHash = string(pwHash)
	user.UpdatedAt = &now
	return svc.userRepo.Update(ctx, user)
}
