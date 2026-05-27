package grpc

import (
	"context"

	"post-service/internal/dto"
	"post-service/internal/services"

	postpb "cosmix/shared/grpc/gen/go/post"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type PostServer struct {
	postpb.UnimplementedPostServiceServer

	postService    *services.PostService
	likeService    *services.LikeService
	commentService *services.CommentService
}

func NewPostServer(
	postService *services.PostService,
	likeService *services.LikeService,
	commentService *services.CommentService,
) *PostServer {
	return &PostServer{
		postService:    postService,
		likeService:    likeService,
		commentService: commentService,
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

	_, err := srv.postService.CreatePost(
		ctx,
		uint(req.AuthorId),
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
	_, err := srv.postService.GetPostByID(
		ctx,
		uint(req.PostId),
	)
	if err != nil {
		return nil, err
	}

	return &postpb.PostResponse{
		Post: nil,
	}, nil
}

func (srv *PostServer) GetFeed(ctx context.Context, req *postpb.GetFeedRequest) (*postpb.PostListResponse, error) {
	result, err := srv.postService.GetFeed(
		ctx,
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
		Id:            uint64(post.ID),
		Content:       post.Content,
		LikesCount:    int32(post.LikesCount),
		CommentsCount: int32(post.CommentsCount),
		CreatedAt:     timestamppb.New(post.CreatedAt),
		User:          mapUser(post.User),
	}

	if post.UpdatedAt != nil {
		result.UpdatedAt =
			timestamppb.New(
				*post.UpdatedAt,
			)
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
	return &postpb.Comment{
		Id:        uint64(comment.ID),
		PostId:    uint64(comment.PostID),
		AuthorId:  uint64(comment.AuthorID),
		Content:   comment.Content,
		CreatedAt: timestamppb.New(comment.CreatedAt),
		UpdatedAt: timestamppb.New(comment.UpdatedAt),
	}
}

func mapMedia(media dto.Media) *postpb.Media {
	result := &postpb.Media{
		Id:        uint64(media.ID),
		PostId:    uint64(media.PostID),
		PublicId:  media.PublicID,
		Url:       media.URL,
		Type:      media.Type,
		CreatedAt: timestamppb.New(media.CreatedAt),
	}

	if media.Duration != nil {
		result.Duration =
			int32(*media.Duration)
	}

	if media.UpdatedAt != nil {
		result.UpdatedAt = timestamppb.New(*media.UpdatedAt)
	}

	return result
}

func mapUser(user dto.User) *postpb.User {
	result := &postpb.User{
		Id:          uint64(user.ID),
		Email:       user.Email,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		CreatedAt:   timestamppb.New(user.CreatedAt),
	}

	if user.UpdatedAt != nil {
		result.UpdatedAt = timestamppb.New(*user.UpdatedAt)
	}

	return result
}
