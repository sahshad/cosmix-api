package app

import (
	"auth-service/internal/controllers"
	"auth-service/internal/messaging"
	"auth-service/internal/repositories"
	"auth-service/internal/services"

	"gorm.io/gorm"
)

type Container struct {
	DB     *gorm.DB
	Rabbit *messaging.Rabbit

	// Controllers
	AuthController *controllers.AuthController

	// Services
	AuthService services.AuthServiceInterface
}

func NewContainer(db *gorm.DB, rabbit *messaging.Rabbit) *Container {

	userRepo := repositories.NewUserRepository(db)
	sessionRepo := repositories.NewAuthSessionRepository(db)

	sessionService := services.NewAuthSessionService(sessionRepo)
	authService := services.NewAuthService(userRepo, sessionService)

	authController := controllers.NewAuthController(authService, rabbit.Channel)

	return &Container{
		DB:             db,
		Rabbit:         rabbit,
		AuthController: authController,
		AuthService:    authService,
	}
}
