package routes

import (
	"auth-service/internal/app"
	"cosmix/shared/core/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, container *app.Container) {
	api := router.Group("/")

	authController := container.AuthController

	api.GET("/health", middleware.ErrorHandler(authController.HealthCheck))
	api.POST("/register", middleware.ErrorHandler(authController.Register))
	api.POST("/login", middleware.ErrorHandler(authController.Login))
	api.GET("/refresh", middleware.ErrorHandler(authController.Refresh))
	api.POST("/logout", middleware.ErrorHandler(authController.Logout))
	api.PUT("/update-password", middleware.ErrorHandler(authController.UpdateUserPassword))
}
