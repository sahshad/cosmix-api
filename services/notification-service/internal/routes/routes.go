package routes

import (
	"notification-service/internal/app"
	"cosmix/shared/core/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, container *app.Container) {
	api := router.Group("/")

	notificationController := container.NotificationController

	api.GET("/health", middleware.ErrorHandler(notificationController.HealthCheck))
	api.GET("/me", middleware.ErrorHandler(notificationController.GetUserNotification))
}