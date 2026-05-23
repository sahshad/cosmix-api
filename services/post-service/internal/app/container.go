package app

import (
	"cosmix/shared/core/rabbitmq"
	"post-service/internal/controllers"
	grpcServer "post-service/internal/grpc/server"
	"post-service/internal/repositories"
	"post-service/internal/services"

	"gorm.io/gorm"
)

type Container struct {
	DB     *gorm.DB
	Rabbit *rabbitmq.Rabbit

	// Controllers
	PostController    *controllers.PostController
	LikeController    *controllers.LikeController
	CommentController *controllers.CommentController

	// Services
	PostService    services.PostServiceInterface
	LikeService    services.LikeServiceInterface
	CommentService services.CommentServiceInterface
	PostUserSvc    services.PostUserServiceInterface

	PostGrpcServer *grpcServer.PostServer
}

func NewContainer(db *gorm.DB, rabbit *rabbitmq.Rabbit) *Container {

	postUserRepo := repositories.NewPostUserRepository(db)
	postUserService := services.NewPostUserService(postUserRepo)

	postRepo := repositories.NewPostRepository(db)
	postService := services.NewPostService(postRepo)
	postController := controllers.NewPostController(postService)

	likeRepo := repositories.NewLikeRepository(db)
	likeService := services.NewLikeService(likeRepo, postRepo)
	likeController := controllers.NewLikeController(likeService)

	commentRepo := repositories.NewCommentRepository(db)
	commentService := services.NewCommentService(commentRepo, postRepo)
	commentController := controllers.NewCommentController(commentService)

	postGrpcServer :=
		grpcServer.NewPostServer(
			postService,
			likeService,
			commentService,
		)

	return &Container{
		DB:     db,
		Rabbit: rabbit,
		// Controllers
		PostController:    postController,
		LikeController:    likeController,
		CommentController: commentController,
		// Services
		PostService:    postService,
		LikeService:    likeService,
		CommentService: commentService,
		PostUserSvc:    postUserService,

		PostGrpcServer: postGrpcServer,
	}
}
