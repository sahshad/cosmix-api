package app

import (
	"auth-service/internal/controllers"
	"auth-service/internal/repositories"
	"auth-service/internal/services"
	"cosmix/shared/core/rabbitmq"

	grpcServer "auth-service/internal/grpc/server"
	"gorm.io/gorm"
)

type Container struct {
	DB     *gorm.DB
	Rabbit *rabbitmq.Rabbit

	// Controllers
	AuthController *controllers.AuthController

	// Services
	AuthService services.AuthServiceInterface

	AuthGrpcServer *grpcServer.AuthServer
}

func NewContainer(db *gorm.DB, rabbit *rabbitmq.Rabbit) *Container {

	userRepo := repositories.NewUserRepository(db)
	sessionRepo := repositories.NewAuthSessionRepository(db)

	sessionService := services.NewAuthSessionService(sessionRepo)
	authService := services.NewAuthService(userRepo, sessionService)

	authController := controllers.NewAuthController(authService, rabbit.Channel)

	authGrpcServer :=
		grpcServer.NewAuthServer(
			authService,
		)

	return &Container{
		DB:             db,
		Rabbit:         rabbit,
		AuthController: authController,
		AuthService:    authService,
		AuthGrpcServer: authGrpcServer,
	}
}
