package routes

import (
	app "auth-service/internal/app"
	"auth-service/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, container *app.Container) {
	api := router.Group("/")

	authController :=  container.AuthController

	api.GET("/health", middlewares.ErrorHandler(authController.HealthCheck))
	api.POST("/register", middlewares.ErrorHandler(authController.Register))
	api.POST("/login", middlewares.ErrorHandler(authController.Login))
	api.GET("/refresh", middlewares.ErrorHandler(authController.Refresh))
	api.POST("/logout", middlewares.ErrorHandler(authController.Logout))
	api.PUT("/update-password", middlewares.ErrorHandler(authController.UpdateUserPassword))
}
