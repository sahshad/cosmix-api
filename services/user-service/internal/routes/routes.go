package routes

import (
	"user-service/internal/app"
	"user-service/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, container *app.Container) {
	api := router.Group("/")

	userProfileController := container.UserProfileController
	followController := container.FollowController

	api.GET("/health", middlewares.ErrorHandler(userProfileController.HealthCheck))
	api.GET("/me", middlewares.ErrorHandler(userProfileController.GetMe))
	api.PUT("/me", middlewares.ErrorHandler(userProfileController.UpdateMe))
	api.GET("/username/:username", middlewares.ErrorHandler(userProfileController.GetByUsername))
	api.POST("/:id/follow", followController.Follow)
	api.DELETE("/:id/unfollow", followController.Unfollow)
	api.GET("/:id/followers", followController.GetFollowers)
	api.GET("/:id/following", followController.GetFollowing)
}
