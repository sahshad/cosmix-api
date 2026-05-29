package services

import (
	"context"
	"errors"
	"time"

	"auth-service/internal/dto"
	"auth-service/internal/models"
	"auth-service/internal/repositories"
	"auth-service/internal/utils"

	"cosmix/shared/core/eventbus"
	"cosmix/shared/core/rabbitmq"

	appErr "cosmix/shared/core/errors"
	authEvents "cosmix/shared/events/auth"

	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type AuthUserService struct {
	authUserRepo   *repositories.AuthUserRepository
	authTokenRepo  *repositories.AuthTokenRepository
	authSessionSvc *AuthSessionService
	rabbitCh       *amqp.Channel
}

func NewAuthUserService(
	authUserRepo *repositories.AuthUserRepository,
	authTokenRepo *repositories.AuthTokenRepository,
	authSessionSvc *AuthSessionService,
	rabbitCh *amqp.Channel,
) *AuthUserService {
	return &AuthUserService{
		authUserRepo:   authUserRepo,
		authTokenRepo:  authTokenRepo,
		authSessionSvc: authSessionSvc,
		rabbitCh:       rabbitCh,
	}
}

func (svc *AuthUserService) Register(ctx context.Context, input dto.RegisterDTO) (*models.AuthUser, error) {
	if _, err := svc.authUserRepo.FindByEmail(ctx, input.Email); err == nil {
		return nil, appErr.NewBadRequest("email", "Email already in use")
	}

	pwHash, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := &models.AuthUser{
		Email:        input.Email,
		DisplayName:  input.DisplayName,
		PasswordHash: string(pwHash),
	}

	if err := svc.authUserRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	tokenValue, err := utils.GenerateSecureToken(32)
	if err != nil {
		return nil, err
	}

	token := &models.AuthToken{
		AuthUserID: user.ID,
		Type:       models.AuthTokenTypeEmailVerification,
		Token:      tokenValue,
		ExpiresAt:  time.Now().Add(24 * time.Hour),
	}

	if err := svc.authTokenRepo.Create(ctx, token); err != nil {
		return nil, err
	}

	requestID := eventbus.RequestID(ctx)
	ctx = eventbus.WithCorrelationID(ctx, requestID)

	event := authEvents.AuthUserEmailVerification{
		EventVersion: authEvents.EventVersionOne,
		Email:        user.Email,
		DisplayName:  input.DisplayName,
		Token:        tokenValue,
	}

	// @Publish - user email verification requested
	if err := eventbus.Publish(
		ctx,
		svc.rabbitCh,
		rabbitmq.ExchangeEvents,
		rabbitmq.AuthUserEmailVerificationRequested,
		event,
	); err != nil {
		return nil, err
	}

	return user, nil
}

func (svc *AuthUserService) VerifyEmail(ctx context.Context, input dto.VerifyEmailDTO) error {
	tokenRecord, err := svc.authTokenRepo.FindByToken(ctx, input.Token)
	if err != nil {
		return appErr.NewBadRequest("token", "Invalid or expired token")
	}

	if time.Now().After(tokenRecord.ExpiresAt) {
		return appErr.NewBadRequest("token", "Token has expired")
	}

	user, err := svc.authUserRepo.FindByID(ctx, tokenRecord.AuthUserID)
	if err != nil {
		return appErr.NewNotFound("user")
	}

	if user.Email != input.Email {
		return appErr.NewBadRequest("email", "Email does not match token")
	}

	if err := utils.VerifyPassword(input.Password, user.PasswordHash); err != nil {
		return appErr.NewBadRequest("password", "Invalid password")
	}

	if user.EmailVerified {
		return appErr.NewBadRequest("email", "Email already verified")
	}

	now := time.Now()
	user.EmailVerified = true
	user.UpdatedAt = &now

	if err := svc.authUserRepo.Update(ctx, user); err != nil {
		return err
	}

	if err := svc.authTokenRepo.Delete(ctx, tokenRecord); err != nil {
		return err
	}

	displayName := user.DisplayName
	username := utils.GenerateUsername(displayName)

	requestID := eventbus.RequestID(ctx)
	ctx = eventbus.WithCorrelationID(ctx, requestID)

	event := authEvents.AuthUserEmailVerificationCompleted{
		EventVersion: authEvents.EventVersionOne,
		AuthUserID:   user.ID,
		Email:        user.Email,
		DisplayName:  displayName,
		Username:     username,
		CreatedAt:    user.CreatedAt,
	}

	// @Publish - user email verification completed
	if err := eventbus.Publish(
		ctx,
		svc.rabbitCh,
		rabbitmq.ExchangeEvents,
		rabbitmq.AuthUserEmailVerificationCompleted,
		event,
	); err != nil {
		return err
	}

	return nil
}

func (svc *AuthUserService) Login(ctx context.Context, input dto.LoginDTO) (*dto.LoginResponseDTO, error) {
	authUser, err := svc.authUserRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, appErr.NewBadRequest("email", "Invalid credentials")
	}

	if !authUser.EmailVerified {
		return nil, appErr.NewBadRequest("email", "Email not verified")
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
	if err := svc.authUserRepo.Update(ctx, authUser); err != nil {
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

	if _, err := svc.authUserRepo.FindByID(ctx, claims.UserID); err != nil {
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
	user, err := svc.authUserRepo.FindByID(ctx, id)
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
	return svc.authUserRepo.Update(ctx, user)
}

func (svc *AuthUserService) ForgotPassword(ctx context.Context, email string) error {
	user, err := svc.authUserRepo.FindByEmail(ctx, email)
	if err != nil {
		return appErr.NewNotFound("user")
	}

	if !user.EmailVerified {
		return appErr.NewBadRequest("email", "Email not verified")
	}

	tokenValue, err := utils.GenerateSecureToken(32)
	if err != nil {
		return err
	}

	token := &models.AuthToken{
		AuthUserID: user.ID,
		Type:       models.AuthTokenTypePasswordReset,
		Token:      tokenValue,
		ExpiresAt:  time.Now().Add(24 * time.Hour),
	}

	if err := svc.authTokenRepo.Create(ctx, token); err != nil {
		return err
	}

	requestID := eventbus.RequestID(ctx)
	ctx = eventbus.WithCorrelationID(ctx, requestID)

	event := authEvents.AuthUserForgotPasswordRequest{
		EventVersion: authEvents.EventVersionOne,
		Email:        user.Email,
		DisplayName:  user.DisplayName,
		Token:        tokenValue,
	}

	// @Publish - user password forgot requested
	if err := eventbus.Publish(
		ctx,
		svc.rabbitCh,
		rabbitmq.ExchangeEvents,
		rabbitmq.AuthUserForgotPasswordRequested,
		event,
	); err != nil {
		return err
	}

	return nil
}

func (svc *AuthUserService) ResetPassword(ctx context.Context, req dto.ResetPasswordDTO) error {
	tokenRecord, err := svc.authTokenRepo.FindByToken(ctx, req.Token)
	if err != nil {
		return appErr.NewBadRequest("token", "Invalid or expired token")
	}

	if time.Now().After(tokenRecord.ExpiresAt) {
		return appErr.NewBadRequest("token", "Token has expired")
	}

	user, err := svc.authUserRepo.FindByID(ctx, tokenRecord.AuthUserID)
	if err != nil {
		return appErr.NewNotFound("user")
	}

	if err := utils.VerifyPassword(req.CurrentPassword, user.PasswordHash); err != nil {
		return appErr.NewBadRequest("current_password", "Incorrect current password")
	}

	pwHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	user.PasswordHash = string(pwHash)
	user.UpdatedAt = &now

	if err := svc.authUserRepo.Update(ctx, user); err != nil {
		return err
	}

	if err := svc.authTokenRepo.Delete(ctx, tokenRecord); err != nil {
		return err
	}

	requestID := eventbus.RequestID(ctx)
	ctx = eventbus.WithCorrelationID(ctx, requestID)

	event := authEvents.AuthUserPasswordChanged{
		EventVersion: authEvents.EventVersionOne,
		AuthUserID:   user.ID,
		Email:        user.Email,
		DisplayName:  user.DisplayName,
		UpdatedAt:    now,
	}

	// @Publish - user password changed
	if err := eventbus.Publish(
		ctx,
		svc.rabbitCh,
		rabbitmq.ExchangeEvents,
		rabbitmq.AuthUserPasswordChanged,
		event,
	); err != nil {
		return err
	}

	return nil
}