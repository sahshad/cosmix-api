package app

import (
	"cosmix/shared/core/rabbitmq"
	// "user-service/internal/controllers"
	grpcServer "user-service/internal/grpc/server"
	"user-service/internal/repositories"
	"user-service/internal/services"

	"gorm.io/gorm"
)

type Container struct {
	DB     *gorm.DB
	Rabbit *rabbitmq.Rabbit

	// Controllers
	// UserProfileController *controllers.UserProfileController
	// FollowController      *controllers.FollowController

	// Services
	UserService   *services.UserService
	FollowService *services.FollowService

	UserGrpcServer *grpcServer.UserServer
}

func NewContainer(db *gorm.DB, rabbit *rabbitmq.Rabbit) *Container {

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	// userController := controllers.NewUserProfileController(userProfileService, rabbit.Channel)

	followRepo := repositories.NewFollowRepository(db)
	followService := services.NewFollowService(followRepo)
	// followController := controllers.NewFollowController(followService)

	userGrpcServer :=
		grpcServer.NewUserServer(
			userService,
			followService,
		)

	return &Container{
		DB:     db,
		Rabbit: rabbit,
		// Controllers
		// UserProfileController: userProfileController,
		// FollowController:      followController,
		// Services
		UserService:   userService,
		FollowService: followService,

		UserGrpcServer: userGrpcServer,
	}
}
