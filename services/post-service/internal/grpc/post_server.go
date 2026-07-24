package grpc

import (
	"context"
	"time"

	"post-service/internal/dto"
	"post-service/internal/services"

	postpb "cosmix/shared/grpc/gen/go/post"

	"github.com/google/uuid"
)

type PostServer struct {
	postpb.UnimplementedPostServiceServer

	postService        *services.PostService
	likeService        *services.LikeService
	commentService     *services.CommentService
	commentLikeService *services.CommentLikeService
}

func NewPostServer(
	postService *services.PostService,
	likeService *services.LikeService,
	commentService *services.CommentService,
	commentLikeService *services.CommentLikeService,
) *PostServer {
	return &PostServer{
		postService:        postService,
		likeService:        likeService,
		commentService:     commentService,
		commentLikeService: commentLikeService,
	}
}

func (srv *PostServer) CreatePost(ctx context.Context, req *postpb.CreatePostRequest) (*postpb.PostResponse, error) {
	input := &dto.CreatePostRequest{
		Content: req.Content,
	}

	for _, media := range req.Media {
		input.Media = append(
			input.Media,
			dto.MediaItem{
				PublicID: media.PublicId,
				URL:      media.Url,
				Type:     media.Type,
				Duration: int(media.Duration),
			},
		)
	}

	authorID, err := uuid.Parse(req.AuthorId)
	if err != nil {
		return nil, err
	}

	_, err = srv.postService.CreatePost(
		ctx,
		authorID,
		input,
	)
	if err != nil {
		return nil, err
	}

	return &postpb.PostResponse{
		Post: nil,
	}, nil
}

func (srv *PostServer) GetPost(ctx context.Context, req *postpb.GetPostRequest) (*postpb.PostResponse, error) {
	postID, err := uuid.Parse(req.PostId)
	if err != nil {
		return nil, err
	}

	_, err = srv.postService.GetPostByID(
		ctx,
		postID,
	)
	if err != nil {
		return nil, err
	}

	return &postpb.PostResponse{
		Post: nil,
	}, nil
}

func (srv *PostServer) DeletePost(ctx context.Context, req *postpb.DeletePostRequest) (*postpb.MessageResponse, error) {
	postID, err := uuid.Parse(req.PostId)
	if err != nil {
		return nil, err
	}

	authorID, err := uuid.Parse(req.AuthorId)
	if err != nil {
		return nil, err
	}

	if err := srv.postService.DeletePost(ctx, postID, authorID); err != nil {
		return nil, err
	}

	return &postpb.MessageResponse{
		Message: "post deleted",
	}, nil
}

func (srv *PostServer) LikePost(ctx context.Context, req *postpb.LikePostRequest) (*postpb.MessageResponse, error) {
	postID, err := uuid.Parse(req.PostId)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	if err := srv.likeService.LikePost(ctx, postID, userID); err != nil {
		return nil, err
	}

	return &postpb.MessageResponse{
		Message: "post liked",
	}, nil
}

func (srv *PostServer) UnlikePost(ctx context.Context, req *postpb.UnlikePostRequest) (*postpb.MessageResponse, error) {
	postID, err := uuid.Parse(req.PostId)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	if err := srv.likeService.UnlikePost(ctx, postID, userID); err != nil {
		return nil, err
	}

	return &postpb.MessageResponse{
		Message: "post unliked",
	}, nil
}

func (srv *PostServer) CreateComment(ctx context.Context, req *postpb.CreateCommentRequest) (*postpb.CommentResponse, error) {
	postID, err := uuid.Parse(req.PostId)
	if err != nil {
		return nil, err
	}

	authorID, err := uuid.Parse(req.AuthorId)
	if err != nil {
		return nil, err
	}

	var parentCommentID *uuid.UUID
	if req.ParentCommentId != "" {
		parsed, err := uuid.Parse(req.ParentCommentId)
		if err != nil {
			return nil, err
		}
		parentCommentID = &parsed
	}

	comment, err := srv.commentService.CreateComment(
		ctx,
		postID,
		authorID,
		&dto.CreateCommentRequest{
			Content:         req.Content,
			ParentCommentID: parentCommentID,
		},
	)
	if err != nil {
		return nil, err
	}

	avatarURL := ""
	if comment.User.AvatarURL != nil {
		avatarURL = *comment.User.AvatarURL
	}

	return &postpb.CommentResponse{
		Comment: mapComment(dto.CommentList{
			ID:              comment.ID,
			PostID:          comment.PostID,
			AuthorID:        comment.UserID,
			ParentCommentID: comment.ParentCommentID,
			Content:         comment.Content,
			LikesCount:      comment.LikesCount,
			RepliesCount:    comment.RepliesCount,
			IsOwner:         true,
			CreatedAt:       comment.CreatedAt.Format(time.RFC3339),
			UpdatedAt:       formatTimePtr(comment.UpdatedAt),
			User: dto.User{
				ID:          comment.User.UserID,
				Username:    comment.User.Username,
				DisplayName: comment.User.DisplayName,
				AvatarURL:   avatarURL,
				CreatedAt:   comment.User.CreatedAt.Format(time.RFC3339),
				UpdatedAt:   formatTimePtr(comment.User.UpdatedAt),
			},
		}),
	}, nil
}

func (srv *PostServer) DeleteComment(ctx context.Context, req *postpb.DeleteCommentRequest) (*postpb.MessageResponse, error) {
	commentID, err := uuid.Parse(req.CommentId)
	if err != nil {
		return nil, err
	}

	authorID, err := uuid.Parse(req.AuthorId)
	if err != nil {
		return nil, err
	}

	if err := srv.commentService.DeleteComment(ctx, commentID, authorID); err != nil {
		return nil, err
	}

	return &postpb.MessageResponse{
		Message: "comment deleted",
	}, nil
}

func (srv *PostServer) LikeComment(ctx context.Context, req *postpb.LikeCommentRequest) (*postpb.MessageResponse, error) {
	commentID, err := uuid.Parse(req.CommentId)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	if err := srv.commentLikeService.LikeComment(ctx, commentID, userID); err != nil {
		return nil, err
	}

	return &postpb.MessageResponse{
		Message: "comment liked",
	}, nil
}

func (srv *PostServer) UnlikeComment(ctx context.Context, req *postpb.UnlikeCommentRequest) (*postpb.MessageResponse, error) {
	commentID, err := uuid.Parse(req.CommentId)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	if err := srv.commentLikeService.UnlikeComment(ctx, commentID, userID); err != nil {
		return nil, err
	}

	return &postpb.MessageResponse{
		Message: "comment unliked",
	}, nil
}

func (srv *PostServer) GetComments(ctx context.Context, req *postpb.GetCommentsRequest) (*postpb.CommentListResponse, error) {
	postID, err := uuid.Parse(req.PostId)
	if err != nil {
		return nil, err
	}

	viewerID, err := parseOptionalUUID(req.UserId)
	if err != nil {
		return nil, err
	}

	result, err := srv.commentService.GetCommentsByPostID(
		ctx,
		postID,
		viewerID,
		&dto.PaginationRequest{
			Page:  req.Page,
			Limit: req.Limit,
		},
	)
	if err != nil {
		return nil, err
	}

	return mapCommentListResponse(result), nil
}

func (srv *PostServer) GetReplies(ctx context.Context, req *postpb.GetRepliesRequest) (*postpb.CommentListResponse, error) {
	commentID, err := uuid.Parse(req.CommentId)
	if err != nil {
		return nil, err
	}

	viewerID, err := parseOptionalUUID(req.UserId)
	if err != nil {
		return nil, err
	}

	result, err := srv.commentService.GetReplies(
		ctx,
		commentID,
		viewerID,
		&dto.PaginationRequest{
			Page:  req.Page,
			Limit: req.Limit,
		},
	)
	if err != nil {
		return nil, err
	}

	return mapCommentListResponse(result), nil
}

func mapCommentListResponse(result *dto.CommentListResponse) *postpb.CommentListResponse {
	response := &postpb.CommentListResponse{
		Pagination: &postpb.Pagination{
			TotalCount: result.Pagination.TotalCount,
			Page:       result.Pagination.Page,
			Limit:      result.Pagination.Limit,
			TotalPages: result.Pagination.TotalPages,
		},
	}

	for _, comment := range result.Comments {
		response.Comments = append(
			response.Comments,
			mapComment(comment),
		)
	}

	return response
}

func parseOptionalUUID(value string) (uuid.UUID, error) {
	if value == "" {
		return uuid.Nil, nil
	}
	return uuid.Parse(value)
}

func (srv *PostServer) GetFeed(ctx context.Context, req *postpb.GetFeedRequest) (*postpb.PostListResponse, error) {
	viewerID, err := parseOptionalUUID(req.UserId)
	if err != nil {
		return nil, err
	}

	result, err := srv.postService.GetFeed(
		ctx,
		viewerID,
		&dto.PaginationRequest{
			Page:  req.Page,
			Limit: req.Limit,
		},
	)
	if err != nil {
		return nil, err
	}

	response := &postpb.PostListResponse{
		Pagination: &postpb.Pagination{
			TotalCount: result.Pagination.TotalCount,
			Page:       result.Pagination.Page,
			Limit:      result.Pagination.Limit,
			TotalPages: result.Pagination.TotalPages,
		},
	}

	for _, post := range result.Posts {
		response.Posts = append(
			response.Posts,
			mapPost(post),
		)
	}

	return response, nil
}

func (srv *PostServer) GetUserPosts(ctx context.Context, req *postpb.GetUserPostsRequest) (*postpb.PostListResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	viewerID, err := parseOptionalUUID(req.ViewerId)
	if err != nil {
		return nil, err
	}

	result, err := srv.postService.GetUserPosts(
		ctx,
		userID,
		viewerID,
		&dto.PaginationRequest{
			Page:  req.Page,
			Limit: req.Limit,
		},
	)
	if err != nil {
		return nil, err
	}

	response := &postpb.PostListResponse{
		Pagination: &postpb.Pagination{
			TotalCount: result.Pagination.TotalCount,
			Page:       result.Pagination.Page,
			Limit:      result.Pagination.Limit,
			TotalPages: result.Pagination.TotalPages,
		},
	}

	for _, post := range result.Posts {
		response.Posts = append(
			response.Posts,
			mapPost(post),
		)
	}

	return response, nil
}

func mapPost(post dto.PostList) *postpb.Post {
	result := &postpb.Post{
		Id:            post.ID.String(),
		Content:       post.Content,
		LikesCount:    int32(post.LikesCount),
		CommentsCount: int32(post.CommentsCount),
		SharesCount:   int32(post.SharesCount),
		IsLiked:       post.IsLiked,
		IsOwner:       post.IsOwner,
		CreatedAt:     post.CreatedAt,
		UpdatedAt:     post.UpdatedAt,
		User:          mapUser(post.User),
	}

	for _, media := range post.Media {
		result.Media = append(
			result.Media,
			mapMedia(media),
		)
	}

	return result
}

func mapComment(comment dto.CommentList) *postpb.Comment {
	parentCommentID := ""
	if comment.ParentCommentID != nil {
		parentCommentID = comment.ParentCommentID.String()
	}

	return &postpb.Comment{
		Id:              comment.ID.String(),
		PostId:          comment.PostID.String(),
		AuthorId:        comment.AuthorID.String(),
		Content:         comment.Content,
		LikesCount:      int32(comment.LikesCount),
		RepliesCount:    int32(comment.RepliesCount),
		IsLiked:         comment.IsLiked,
		IsOwner:         comment.IsOwner,
		ParentCommentId: parentCommentID,
		CreatedAt:       comment.CreatedAt,
		UpdatedAt:       comment.UpdatedAt,
		User:            mapUser(comment.User),
	}
}

func mapMedia(media dto.Media) *postpb.Media {
	result := &postpb.Media{
		Id:        media.ID.String(),
		PostId:    media.PostID.String(),
		PublicId:  media.PublicID,
		Url:       media.URL,
		Type:      media.Type,
		CreatedAt: media.CreatedAt,
		UpdatedAt: media.UpdatedAt,
	}

	if media.Duration != nil {
		result.Duration =
			int32(*media.Duration)
	}

	return result
}

func mapUser(user dto.User) *postpb.User {
	return &postpb.User{
		Id:          user.ID.String(),
		Email:       user.Email,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		AvatarUrl:   user.AvatarURL,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

func formatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format(time.RFC3339)
}
