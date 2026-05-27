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
	userService := services.NewUserService(userRepo)
	// userController := controllers.NewUserProfileController(userProfileService, rabbit.Channel)

	followRepo := repositories.NewFollowRepository(db)
	followService := services.NewFollowService(followRepo)
	// followController := controllers.NewFollowController(followService)

	userGrpcServer :=
		grpc.NewUserServer(
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
