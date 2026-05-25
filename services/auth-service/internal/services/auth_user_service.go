package services

import (
	"context"
	"errors"
	"time"

	"auth-service/internal/dto"
	"auth-service/internal/models"
	"auth-service/internal/repositories"
	"auth-service/internal/utils"

	// publisher "auth-service/internal/messaging/publisher"
	appErr "cosmix/shared/core/errors"
	"cosmix/shared/core/eventbus"
	"cosmix/shared/core/rabbitmq"
	authEvents "cosmix/shared/events/auth"

	amqp "github.com/rabbitmq/amqp091-go"

	"gorm.io/gorm"
)

// type AuthServiceInterface interface {
// 	Register(ctx context.Context, input dto.RegisterDTO) (*models.AuthUser, error)
// 	Login(ctx context.Context, input dto.LoginDTO) (*dto.LoginResponseDTO, error)
// 	Refresh(ctx context.Context, refreshToken string) (*dto.RefreshResponseDTO, error)
// 	GetByID(ctx context.Context, id uint) (*models.AuthUser, error)
// 	UpdateUserPassword(ctx context.Context, userID uint, newPassword string) error
// }

type AuthUserService struct {
	repo           *repositories.AuthUserRepository
	authSessionSvc *AuthSessionService
	rabbitCh       *amqp.Channel
}

func NewAuthUserService(
	repo *repositories.AuthUserRepository,
	authSessionSvc *AuthSessionService,
	rabbitCh *amqp.Channel,
) *AuthUserService {
	return &AuthUserService{
		repo:           repo,
		authSessionSvc: authSessionSvc,
		rabbitCh:       rabbitCh,
	}
}

func (svc *AuthUserService) Register(ctx context.Context, input dto.RegisterDTO) (*models.AuthUser, error) {
	if _, err := svc.repo.FindByEmail(ctx, input.Email); err == nil {
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

	if err := svc.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	username := utils.GenerateUsername(input.DisplayName)

	event := authEvents.AuthUserRegistered{
		EventVersion: authEvents.EventVersionOne,
		AuthUserID:   user.ID,
		Email:        user.Email,
		Username:     username,
		DisplayName:  input.DisplayName,
		CreatedAt:    time.Now().UTC(),
	}

	requestID := eventbus.RequestID(ctx)
	eventbus.WithCorrelationID(ctx, requestID)

	// Publish user registered event
	if err := eventbus.Publish(
		ctx,
		svc.rabbitCh,
		rabbitmq.ExchangeEvents,
		rabbitmq.AuthUserRegistered,
		event,
	); err != nil {
		// handle error
	}
	// publisher.PublishAuthUserRegistered(svc.rabbitCh, authEvents.AuthUserRegistered{
	// 	EventVersion: authEvents.EventVersionOne,
	// 	AuthUserID:   user.ID,
	// 	Email:        user.Email,
	// 	Username:     username,
	// 	DisplayName:  input.DisplayName,
	// 	CreatedAt:    time.Now().UTC(),
	// })

	return user, nil
}

func (svc *AuthUserService) Login(ctx context.Context, input dto.LoginDTO) (*dto.LoginResponseDTO, error) {
	authUser, err := svc.repo.FindByEmail(ctx, input.Email)
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
	if err := svc.repo.Update(ctx, authUser); err != nil {
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

	if _, err := svc.authSessionSvc.Create(ctx, &session); err != nil {
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

func (svc *AuthUserService) Refresh(ctx context.Context, refreshToken string) (*dto.RefreshResponseDTO, error) {
	claims, err := utils.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, appErr.NewUnauthorized("Invalid Token")
	}

	if _, err := svc.repo.FindByID(ctx, claims.UserID); err != nil {
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

func (svc *AuthUserService) GetByID(ctx context.Context, id uint) (*models.AuthUser, error) {
	user, err := svc.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErr.NewNotFound("user")
		}
		return nil, err
	}
	return user, nil
}

func (svc *AuthUserService) UpdateUserPassword(ctx context.Context, userID uint, newPassword string) error {
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
	return svc.repo.Update(ctx, user)
}
