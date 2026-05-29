package app

import (
	"auth-service/internal/grpc"
	"auth-service/internal/repositories"
	"auth-service/internal/services"

	"cosmix/shared/core/rabbitmq"

	"gorm.io/gorm"
)

type Container struct {
	DB     *gorm.DB
	Rabbit *rabbitmq.Rabbit

	AuthUserService *services.AuthUserService

	AuthGrpcServer *grpc.AuthServer
}

func NewContainer(db *gorm.DB, rabbit *rabbitmq.Rabbit) *Container {

	authUserRepo := repositories.NewAuthUserRepository(db)
	authSessionRepo := repositories.NewAuthSessionRepository(db)
	authTokenRepo := repositories.NewAuthTokenRepository(db)

	authSessionService := services.NewAuthSessionService(authSessionRepo)
	authUserService := services.NewAuthUserService(authUserRepo, authTokenRepo, authSessionService, rabbit.Channel)

	authGrpcServer := grpc.NewAuthServer(authUserService)

	return &Container{
		DB:     db,
		Rabbit: rabbit,

		AuthUserService: authUserService,

		AuthGrpcServer: authGrpcServer,
	}
}
