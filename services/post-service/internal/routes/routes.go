package routes

import (
	"post-service/internal/app"
	"cosmix/shared/core/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, container *app.Container) {
	api := router.Group("/")

	postController := container.PostController
	likeController := container.LikeController
	commentController := container.CommentController

	// Posts
	api.POST("/", middleware.ErrorHandler(postController.CreatePost))
	api.GET("/", middleware.ErrorHandler(postController.GetFeed))
	api.GET("/:id", middleware.ErrorHandler(postController.GetPost))
	api.PUT("/:id", middleware.ErrorHandler(postController.UpdatePost))
	api.DELETE("/:id", middleware.ErrorHandler(postController.DeletePost))

	// Likes
	api.POST("/:id/like", middleware.ErrorHandler(likeController.LikePost))
	api.DELETE("/:id/like", middleware.ErrorHandler(likeController.UnlikePost))

	// Comments
	api.POST("/:id/comment", middleware.ErrorHandler(commentController.CreateComment))
	api.GET("/:id/comment", middleware.ErrorHandler(commentController.GetComments))
	api.PUT("/comment/:commentId", middleware.ErrorHandler(commentController.UpdateComment))
	api.DELETE("/comment/:commentId", middleware.ErrorHandler(commentController.DeleteComment))
}
