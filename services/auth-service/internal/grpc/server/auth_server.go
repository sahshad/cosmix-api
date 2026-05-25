package server

import (
	"context"

	authpb "cosmix/shared/grpc/gen/go/auth"

	"auth-service/internal/dto"
	"auth-service/internal/services"
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

func (s *AuthServer) Register(
	ctx context.Context,
	req *authpb.RegisterRequest,
) (*authpb.RegisterResponse, error) {

	user, err := s.authService.Register(
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

func (s *AuthServer) Login(
	ctx context.Context,
	req *authpb.LoginRequest,
) (*authpb.LoginResponse, error) {

	result, err := s.authService.Login(
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

func (s *AuthServer) Refresh(
	ctx context.Context,
	req *authpb.RefreshRequest,
) (*authpb.RefreshResponse, error) {

	result, err := s.authService.Refresh(
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

func (s *AuthServer) UpdateUserPassword(
	ctx context.Context,
	req *authpb.UpdateUserPasswordRequest,
) (*authpb.UpdateUserPasswordResponse, error) {

	err := s.authService.UpdateUserPassword(
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