package grpc

import (
	"context"

	"auth-service/internal/dto"
	"auth-service/internal/services"

	authpb "cosmix/shared/grpc/gen/go/auth"
)

type AuthServer struct {
	authpb.UnimplementedAuthServiceServer
	authService *services.AuthUserService
}

func NewAuthServer(
	authService *services.AuthUserService,
) *AuthServer {
	return &AuthServer{
		authService: authService,
	}
}

func (srv *AuthServer) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	user, err := srv.authService.Register(
		ctx,
		dto.RegisterDTO{
			DisplayName: req.DisplayName,
			Email:       req.Email,
			Password:    req.Password,
		},
	)
	if err != nil {
		return nil, err
	}

	return &authpb.RegisterResponse{
		UserId: uint64(user.ID),
	}, nil
}

func (srv *AuthServer) VerifyEmail(ctx context.Context, req *authpb.VerifyEmailRequest) (*authpb.VerifyEmailResponse, error) {
	err := srv.authService.VerifyEmail(
		ctx,
		dto.VerifyEmailDTO{
			Token:    req.Token,
			Email:    req.Email,
			Password: req.Password,
		},
	)
	if err != nil {
		return nil, err
	}

	return &authpb.VerifyEmailResponse{
		Message: "email verified successfully",
	}, nil
}

func (srv *AuthServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	result, err := srv.authService.Login(
		ctx,
		dto.LoginDTO{
			Email:    req.Email,
			Password: req.Password,
		},
	)
	if err != nil {
		return nil, err
	}

	resp := &authpb.LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}

	if result.AuthUser != nil {
		resp.User = &authpb.AuthUser{
			Email:         result.AuthUser.Email,
			IsActive:      result.AuthUser.IsActive,
			EmailVerified: result.AuthUser.EmailVerified,
		}
	}

	return resp, nil
}

func (srv *AuthServer) Refresh(ctx context.Context, req *authpb.RefreshRequest) (*authpb.RefreshResponse, error) {
	result, err := srv.authService.Refresh(
		ctx,
		req.RefreshToken,
	)
	if err != nil {
		return nil, err
	}

	return &authpb.RefreshResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}, nil
}

func (srv *AuthServer) UpdateUserPassword(ctx context.Context, req *authpb.UpdateUserPasswordRequest) (*authpb.UpdateUserPasswordResponse, error) {
	err := srv.authService.UpdateUserPassword(
		ctx,
		uint(req.UserId),
		req.NewPassword,
	)
	if err != nil {
		return nil, err
	}

	return &authpb.UpdateUserPasswordResponse{
		Message: "password updated successfully",
	}, nil
}

func (srv *AuthServer) ForgotPassword(ctx context.Context, req *authpb.ForgotPasswordRequest) (*authpb.ForgotPasswordResponse, error) {
	err := srv.authService.ForgotPassword(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	return &authpb.ForgotPasswordResponse{
		Message: "if an account with this email exists, you will receive a reset link shortly",
	}, nil
}

func (srv *AuthServer) ResetPassword(ctx context.Context, req *authpb.ResetPasswordRequest) (*authpb.ResetPasswordResponse, error) {
	err := srv.authService.ResetPassword(ctx, dto.ResetPasswordDTO{
		Token:           req.Token,
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	})
	if err != nil {
		return nil, err
	}

	return &authpb.ResetPasswordResponse{
		Message: "password reset successfully",
	}, nil
}
