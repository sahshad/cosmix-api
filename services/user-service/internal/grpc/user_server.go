package grpc

import (
	"context"

	"user-service/internal/dto"
	"user-service/internal/services"

	userpb "cosmix/shared/grpc/gen/go/user"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserServer struct {
	userpb.UnimplementedUserServiceServer

	userService   *services.UserService
	followService *services.FollowService
}

func NewUserServer(
	userService *services.UserService,
	followService *services.FollowService,
) *UserServer {
	return &UserServer{
		userService:   userService,
		followService: followService,
	}
}

func (srv *UserServer) GetProfile(ctx context.Context, req *userpb.GetProfileRequest) (*userpb.UserProfileResponse, error) {
	profile, err := srv.userService.GetProfile(ctx, uint(req.UserId))
	if err != nil {
		return nil, err
	}

	return &userpb.UserProfileResponse{
		User: mapUser(profile.User),
	}, nil
}

func (srv *UserServer) GetProfileByUsername(ctx context.Context, req *userpb.GetProfileByUsernameRequest) (*userpb.UserProfileResponse, error) {
	profile, err := srv.userService.GetProfileByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	return &userpb.UserProfileResponse{
		User: mapUser(profile.User),
	}, nil
}

func (srv *UserServer) UpdateProfile(ctx context.Context, req *userpb.UpdateProfileRequest) (*userpb.UserProfileResponse, error) {
	profile, err := srv.userService.UpdateProfile(
		ctx,
		uint(req.UserId),
		dto.UpdateProfileDTO{
			DisplayName: req.DisplayName,
			Username:    req.Username,
			DateOfBirth: req.DateOfBirth,
			AvatarURL:   req.AvatarUrl,
			Bio:         req.Bio,
		},
	)
	if err != nil {
		return nil, err
	}

	return &userpb.UserProfileResponse{
		User: mapUser(profile.User),
	}, nil
}

func (srv *UserServer) Follow(ctx context.Context, req *userpb.FollowRequest) (*userpb.FollowResponse, error) {
	err := srv.followService.Follow(
		ctx,
		uint(req.FollowerId),
		uint(req.FollowingId),
	)
	if err != nil {
		return nil, err
	}

	return &userpb.FollowResponse{
		Message: "followed successfully",
	}, nil
}

func (srv *UserServer) Unfollow(ctx context.Context, req *userpb.UnfollowRequest) (*userpb.UnfollowResponse, error) {
	err := srv.followService.Unfollow(
		ctx,
		uint(req.FollowerId),
		uint(req.FollowingId),
	)
	if err != nil {
		return nil, err
	}

	return &userpb.UnfollowResponse{
		Message: "unfollowed successfully",
	}, nil
}

// func (s *UserServer) GetFollowers(
// 	ctx context.Context,
// 	req *userpb.GetFollowersRequest,
// ) (*userpb.UserListResponse, error) {

// 	users, err := s.followService.GetFollowers(
// 		ctx,
// 		uint(req.UserId),
// 	)

// 	if err != nil {
// 		return nil, err
// 	}

// 	response := &userpb.UserListResponse{}

// 	for _, user := range users {
// 		response.Users = append(
// 			response.Users,
// 			mapUser(user),
// 		)
// 	}

// 	return response, nil
// }

// func (s *UserServer) GetFollowing(
// 	ctx context.Context,
// 	req *userpb.GetFollowingRequest,
// ) (*userpb.UserListResponse, error) {

// 	users, err := s.followService.GetFollowing(
// 		ctx,
// 		uint(req.UserId),
// 	)

// 	if err != nil {
// 		return nil, err
// 	}

// 	response := &userpb.UserListResponse{}

// 	for _, user := range users {
// 		response.Users = append(
// 			response.Users,
// 			mapUser(user),
// 		)
// 	}

// 	return response, nil
// }

func mapUser(user dto.UserResponse) *userpb.User {
	var dateOfBirth *timestamppb.Timestamp
	if user.DateOfBirth != nil {
		dateOfBirth = timestamppb.New(
			*user.DateOfBirth,
		)
	}

	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt != nil {
		updatedAt = timestamppb.New(
			*user.UpdatedAt,
		)
	}

	return &userpb.User{
		Id:          uint64(user.UserID),
		DisplayName: user.DisplayName,
		Username:    user.Username,
		Email:       user.Email,
		IsPrivate:   user.IsPrivate,
		IsActive:    user.IsActive,
		DateOfBirth: dateOfBirth,
		AvatarUrl:   user.AvatarURL,
		Bio:         user.Bio,
		CreatedAt:   timestamppb.New(user.CreatedAt),
		UpdatedAt:   updatedAt,
	}
}
