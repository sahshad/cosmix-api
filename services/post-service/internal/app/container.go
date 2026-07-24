package app

import (
	"post-service/internal/grpc"
	"post-service/internal/repositories"
	"post-service/internal/services"

	"cosmix/shared/core/rabbitmq"

	"gorm.io/gorm"
)

type Container struct {
	DB     *gorm.DB
	Rabbit *rabbitmq.Rabbit

	PostService        *services.PostService
	LikeService        *services.LikeService
	CommentService     *services.CommentService
	CommentLikeService *services.CommentLikeService
	PostUserSvc        *services.PostUserService

	PostGrpcServer *grpc.PostServer
}

func NewContainer(db *gorm.DB, rabbit *rabbitmq.Rabbit) *Container {

	postUserRepo := repositories.NewPostUserRepository(db)
	postRepo := repositories.NewPostRepository(db)
	likeRepo := repositories.NewLikeRepository(db)
	commentRepo := repositories.NewCommentRepository(db)
	commentLikeRepo := repositories.NewCommentLikeRepository(db)

	postUserService := services.NewPostUserService(postUserRepo)
	postService := services.NewPostService(postRepo)
	likeService := services.NewLikeService(likeRepo, postRepo)
	commentService := services.NewCommentService(commentRepo, postRepo, postUserRepo)
	commentLikeService := services.NewCommentLikeService(commentLikeRepo, commentRepo)

	postGrpcServer := grpc.NewPostServer(
		postService,
		likeService,
		commentService,
		commentLikeService,
	)

	return &Container{
		DB:     db,
		Rabbit: rabbit,

		PostService:        postService,
		LikeService:        likeService,
		CommentService:     commentService,
		CommentLikeService: commentLikeService,
		PostUserSvc:        postUserService,

		PostGrpcServer: postGrpcServer,
	}
}
