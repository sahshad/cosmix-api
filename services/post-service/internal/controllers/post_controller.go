package controllers

import (
	"cosmix/shared/core/httpx"
	"post-service/internal/dto"
	"post-service/internal/services"

	"github.com/gin-gonic/gin"
)

type PostController struct {
	svc services.PostServiceInterface
}

func NewPostController(svc services.PostServiceInterface) *PostController {
	return &PostController{svc: svc}
}

func (ctrl *PostController) CreatePost(c *gin.Context) (interface{}, error) {
	authorID, err := httpx.ParseUserIDHeader(c)
	if err != nil {
		return nil, err
	}

	var req dto.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	ctx := c.Request.Context()
	post, err := ctrl.svc.CreatePost(ctx, uint(authorID), &req)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (ctrl *PostController) GetPost(c *gin.Context) (interface{}, error) {
	id, err := httpx.ParseParamID(c, "id")
	if err != nil {
		return nil, err
	}

	ctx := c.Request.Context()
	post, err := ctrl.svc.GetPostByID(ctx, uint(id))
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (ctrl *PostController) GetFeed(c *gin.Context) (interface{}, error) {
	limit, err := httpx.ParseQueryIntWithDefault(c, "limit", 20)
	if err != nil {
		return nil, err
	}
	page, err := httpx.ParseQueryIntWithDefault(c, "page", 0)
	if err != nil {
		return nil, err
	}

	params := &dto.PaginationRequest{
		Page:  page,
		Limit: limit,
	}

	ctx := c.Request.Context()
	posts, err := ctrl.svc.GetFeed(ctx, params)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (ctrl *PostController) UpdatePost(c *gin.Context) (interface{}, error) {
	authorID, err := httpx.ParseUserIDHeader(c)
	if err != nil {
		return nil, err
	}

	id, err := httpx.ParseParamID(c, "id")
	if err != nil {
		return nil, err
	}

	var req dto.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	ctx := c.Request.Context()
	post, err := ctrl.svc.UpdatePost(ctx, uint(id), uint(authorID), &req)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (ctrl *PostController) DeletePost(c *gin.Context) (interface{}, error) {
	authorID, err := httpx.ParseUserIDHeader(c)
	if err != nil {
		return nil, err
	}

	id, err := httpx.ParseParamID(c, "id")
	if err != nil {
		return nil, err
	}

	ctx := c.Request.Context()
	if err := ctrl.svc.DeletePost(ctx, uint(id), uint(authorID)); err != nil {
		return nil, err
	}

	return map[string]string{"message": "post deleted successfully"}, nil
}
