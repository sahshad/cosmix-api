package app

import (
	// "auth-service/internal/controllers"
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
	// AuthController *controllers.AuthController

	// Services
	AuthUserService *services.AuthUserService

	AuthGrpcServer *grpcServer.AuthServer
}

func NewContainer(db *gorm.DB, rabbit *rabbitmq.Rabbit) *Container {

	userRepo := repositories.NewAuthUserRepository(db)
	sessionRepo := repositories.NewAuthSessionRepository(db)

	sessionService := services.NewAuthSessionService(sessionRepo)
	authUserService := services.NewAuthUserService(userRepo, sessionService, rabbit.Channel)

	// authController := controllers.NewAuthController(authService, rabbit.Channel)

	authGrpcServer :=
		grpcServer.NewAuthServer(
			authUserService,
		)

	return &Container{
		DB:             db,
		Rabbit:         rabbit,
		// AuthController: authController,
		AuthUserService:    authUserService,
		AuthGrpcServer: authGrpcServer,
	}
}
