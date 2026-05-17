package controllers

import (
	"cosmix/shared/core/httpx"
	"post-service/internal/dto"
	"post-service/internal/services"

	"github.com/gin-gonic/gin"
)

type CommentController struct {
	svc services.CommentServiceInterface
}

func NewCommentController(svc services.CommentServiceInterface) *CommentController {
	return &CommentController{svc: svc}
}

func (ctrl *CommentController) CreateComment(c *gin.Context) (interface{}, error) {
	authorID, err := httpx.ParseUserIDHeader(c)
	if err != nil {
		return nil, err
	}

	postID, err := httpx.ParseParamID(c, "id")
	if err != nil {
		return nil, err
	}

	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	ctx := c.Request.Context()
	comment, err := ctrl.svc.CreateComment(ctx, uint(postID), uint(authorID), &req)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (ctrl *CommentController) GetComments(c *gin.Context) (interface{}, error) {
	postID, err := httpx.ParseParamID(c, "id")
	if err != nil {
		return nil, err
	}

	limit, err := httpx.ParseQueryIntWithDefault(c, "limit", 10)
	if err != nil {
		return nil, err
	}

	page, err := httpx.ParseQueryIntWithDefault(c, "page", 1)
	if err != nil {
		return nil, err
	}

	params := &dto.PaginationRequest{
		Limit: limit,
		Page:  page,
	}

	ctx := c.Request.Context()
	comments, err := ctrl.svc.GetCommentsByPostID(ctx, uint(postID), params)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (ctrl *CommentController) UpdateComment(c *gin.Context) (interface{}, error) {
	authorID, err := httpx.ParseUserIDHeader(c)
	if err != nil {
		return nil, err
	}

	commentID, err := httpx.ParseParamID(c, "commentId")
	if err != nil {
		return nil, err
	}

	var req dto.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, err
	}

	ctx := c.Request.Context()
	comment, err := ctrl.svc.UpdateComment(ctx, uint(commentID), uint(authorID), &req)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (ctrl *CommentController) DeleteComment(c *gin.Context) (interface{}, error) {
	authorID, err := httpx.ParseUserIDHeader(c)
	if err != nil {
		return nil, err
	}

	commentID, err := httpx.ParseParamID(c, "commentId")
	if err != nil {
		return nil, err
	}

	ctx := c.Request.Context()
	if err := ctrl.svc.DeleteComment(ctx, uint(commentID), uint(authorID)); err != nil {
		return nil, err
	}

	return map[string]string{"message": "comment deleted successfully"}, nil
}
