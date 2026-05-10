package app

import (
	"auth-service/internal/controllers"
	"auth-service/internal/repositories"
	"auth-service/internal/services"
	"cosmix/shared/core/rabbitmq"

	"gorm.io/gorm"
)

type Container struct {
	DB     *gorm.DB
	Rabbit *rabbitmq.Rabbit

	// Controllers
	AuthController *controllers.AuthController

	// Services
	AuthService services.AuthServiceInterface
}

func NewContainer(db *gorm.DB, rabbit *rabbitmq.Rabbit) *Container {

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
