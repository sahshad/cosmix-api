package grpc

import (
	"context"

	"user-service/internal/dto"
	"user-service/internal/services"

	appErr "cosmix/shared/core/errors"
	userpb "cosmix/shared/grpc/gen/go/user"

	"github.com/google/uuid"
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

	viewerID, err := parseOptionalUUID(req.ViewerId)
	if err != nil {
		return nil, err
	}

	profile, err := srv.userService.GetProfile(ctx, userId, viewerID)
	if err != nil {
		return nil, err
	}

	return &userpb.UserProfileResponse{
		User: mapUser(profile.User),
	}, nil
}

func (srv *UserServer) GetProfileByUsername(ctx context.Context, req *userpb.GetProfileByUsernameRequest) (*userpb.UserProfileResponse, error) {
	viewerID, err := parseOptionalUUID(req.ViewerId)
	if err != nil {
		return nil, err
	}

	profile, err := srv.userService.GetProfileByUsername(ctx, req.Username, viewerID)
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
			DisplayName:   req.DisplayName,
			Username:      req.Username,
			DateOfBirth:   req.DateOfBirth,
			AvatarURL:     req.AvatarUrl,
			CoverImageURL: req.CoverImageUrl,
			Website:       req.Website,
			Location:      req.Location,
			Bio:           req.Bio,
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

func (srv *UserServer) GetFollowers(
	ctx context.Context,
	req *userpb.GetFollowersRequest,
) (*userpb.UserListResponse, error) {
	return srv.listFollowUsers(ctx, req.UserId, srv.followService.GetFollowers)
}

func (srv *UserServer) GetFollowing(
	ctx context.Context,
	req *userpb.GetFollowingRequest,
) (*userpb.UserListResponse, error) {
	return srv.listFollowUsers(ctx, req.UserId, srv.followService.GetFollowing)
}

func (srv *UserServer) listFollowUsers(
	ctx context.Context,
	rawUserID string,
	load func(context.Context, uuid.UUID) ([]uuid.UUID, error),
) (*userpb.UserListResponse, error) {
	userID, err := uuid.Parse(rawUserID)
	if err != nil {
		return nil, appErr.NewBadRequest("userId", "invalid user id")
	}

	userIDs, err := load(ctx, userID)
	if err != nil {
		return nil, err
	}

	usersByID, err := srv.userService.GetUserResponsesByIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	response := &userpb.UserListResponse{Users: make([]*userpb.User, 0, len(userIDs))}
	for _, id := range userIDs {
		if user, ok := usersByID[id]; ok {
			response.Users = append(response.Users, mapUser(user))
		}
	}

	return response, nil
}

func mapUser(user dto.UserResponse) *userpb.User {
	return &userpb.User{
		Id:             user.UserID.String(),
		DisplayName:    user.DisplayName,
		Username:       user.Username,
		Email:          user.Email,
		IsPrivate:      user.IsPrivate,
		IsVerified:     user.IsVerified,
		IsActive:       user.IsActive,
		DateOfBirth:    user.DateOfBirth,
		AvatarUrl:      user.AvatarURL,
		CoverImageUrl:  user.CoverImageURL,
		Website:        user.Website,
		Location:       user.Location,
		Bio:            user.Bio,
		LastSeenAt:     user.LastSeenAt,
		FollowersCount: int32(user.FollowersCount),
		FollowingCount: int32(user.FollowingCount),
		IsFollowing:    user.IsFollowing,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}

func parseOptionalUUID(value string) (uuid.UUID, error) {
	if value == "" {
		return uuid.Nil, nil
	}
	return uuid.Parse(value)
}
