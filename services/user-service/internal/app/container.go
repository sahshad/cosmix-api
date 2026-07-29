package app

import (
	"user-service/internal/grpc"
	"user-service/internal/repositories"
	"user-service/internal/services"

	"cosmix/shared/core/rabbitmq"

	"gorm.io/gorm"
)

type Container struct {
	DB     *gorm.DB
	Rabbit *rabbitmq.Rabbit

	UserService   *services.UserService
	FollowService *services.FollowService

	UserGrpcServer *grpc.UserServer
}

func NewContainer(db *gorm.DB, rabbit *rabbitmq.Rabbit) *Container {

	userRepo := repositories.NewUserRepository(db)
	followRepo := repositories.NewFollowRepository(db)
	userService := services.NewUserService(userRepo, followRepo, rabbit.Channel)
	followService := services.NewFollowService(followRepo)

	userGrpcServer :=
		grpc.NewUserServer(
			userService,
			followService,
		)

	return &Container{
		DB:     db,
		Rabbit: rabbit,

		UserService:   userService,
		FollowService: followService,

		UserGrpcServer: userGrpcServer,
	}
}
