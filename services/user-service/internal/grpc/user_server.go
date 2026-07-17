package grpc

import (
	"context"

	"user-service/internal/dto"
	"user-service/internal/services"

	userpb "cosmix/shared/grpc/gen/go/user"

	"github.com/google/uuid"
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
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	profile, err := srv.userService.GetProfile(ctx, userId)
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
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	profile, err := srv.userService.UpdateProfile(
		ctx,
		userId,
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
	followerID, err := uuid.Parse(req.FollowerId)
	if err != nil {
		return nil, err
	}

	followingID, err := uuid.Parse(req.FollowingId)
	if err != nil {
		return nil, err
	}

	err = srv.followService.Follow(
		ctx,
		followerID,
		followingID,
	)
	if err != nil {
		return nil, err
	}

	return &userpb.FollowResponse{
		Message: "followed successfully",
	}, nil
}

func (srv *UserServer) Unfollow(ctx context.Context, req *userpb.UnfollowRequest) (*userpb.UnfollowResponse, error) {
	followerID, err := uuid.Parse(req.FollowerId)
	if err != nil {
		return nil, err
	}

	followingID, err := uuid.Parse(req.FollowingId)
	if err != nil {
		return nil, err
	}

	err = srv.followService.Unfollow(
		ctx,
		followerID,
		followingID,
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
		Id:          user.UserID.String(),
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
