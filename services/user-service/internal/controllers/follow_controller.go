package controllers

import (
	"user-service/internal/services"
	"cosmix/shared/core/httpx"

	"github.com/gin-gonic/gin"
)

type FollowController struct {
	service services.FollowServiceInterface
}

func NewFollowController(service services.FollowServiceInterface) *FollowController {
	return &FollowController{service: service}
}

func (ctrl *FollowController) Follow(c *gin.Context) (interface{}, error) {
	followerID, err := httpx.ParseUserIDHeader(c)
	if err != nil {
		return nil, err
	}

	followingID, err := httpx.ParseParamID(c, "id")
	if err != nil {
		return nil, err
	}

	ctx := c.Request.Context()
	
	err = ctrl.service.Follow(ctx, uint(followerID), uint(followingID))
	if err != nil {
		return nil, err
	}

	return gin.H{"message": "followed successfully"}, nil
}

func (ctrl *FollowController) Unfollow(c *gin.Context) (interface{}, error) {
	followerID, err := httpx.ParseUserIDHeader(c)
	if err != nil {
		return nil, err
	}

	followingID, err := httpx.ParseParamID(c, "id")
	if err != nil {
		return nil, err
	}

	ctx := c.Request.Context()
	
	err = ctrl.service.Unfollow(ctx, uint(followerID), uint(followingID))
	if err != nil {
		return nil, err
	}

	return gin.H{"message": "unfollowed successfully"}, nil
}

func (ctrl *FollowController) GetFollowers(c *gin.Context) (interface{}, error) {
	id, err := httpx.ParseParamID(c, "id")
	if err != nil {
		return nil, err
	}

	ctx := c.Request.Context()

	followers, err := ctrl.service.GetFollowers(ctx, uint(id))
	if err != nil {
		return nil, err
	}

	return followers, nil
}

func (ctrl *FollowController) GetFollowing(c *gin.Context) (interface{}, error) {
	id, err := httpx.ParseParamID(c, "id")
	if err != nil {
		return nil, err
	}

	ctx := c.Request.Context()

	following, err := ctrl.service.GetFollowing(ctx, uint(id))
	if err != nil {
		return nil, err
	}

	return following, nil
}
