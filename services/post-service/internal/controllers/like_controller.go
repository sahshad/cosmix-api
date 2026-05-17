package controllers

import (
	"cosmix/shared/core/httpx"
	"post-service/internal/services"

	"github.com/gin-gonic/gin"
)

type LikeController struct {
	svc services.LikeServiceInterface
}

func NewLikeController(svc services.LikeServiceInterface) *LikeController {
	return &LikeController{svc: svc}
}

func (ctrl *LikeController) LikePost(c *gin.Context) (interface{}, error) {
	userID, err := httpx.ParseUserIDHeader(c)
	if err != nil {
		return nil, err
	}

	postID, err := httpx.ParseParamID(c, "id")
	if err != nil {
		return nil, err
	}

	ctx := c.Request.Context()
	if err := ctrl.svc.LikePost(ctx, uint(postID), uint(userID)); err != nil {
		return nil, err
	}

	return map[string]string{"message": "post liked successfully"}, nil
}

func (ctrl *LikeController) UnlikePost(c *gin.Context) (interface{}, error) {
	userID, err := httpx.ParseUserIDHeader(c)
	if err != nil {
		return nil, err
	}

	postID, err := httpx.ParseParamID(c, "id")
	if err != nil {
		return nil, err
	}

	ctx := c.Request.Context()
	if err := ctrl.svc.UnlikePost(ctx, uint(postID), uint(userID)); err != nil {
		return nil, err
	}

	return map[string]string{"message": "post unliked successfully"}, nil
}
