package app

import (
	"user-service/internal/controllers"
	"cosmix/shared/core/rabbitmq"
	"user-service/internal/repositories"
	"user-service/internal/services"

	"gorm.io/gorm"
)

type Container struct {
	DB     *gorm.DB
	Rabbit *rabbitmq.Rabbit

	// Controllers
	UserProfileController *controllers.UserProfileController
	FollowController      *controllers.FollowController

	// Services
	UserProfileService services.UserProfileServiceInterface
	FollowService      services.FollowServiceInterface
}

func NewContainer(db *gorm.DB, rabbit *rabbitmq.Rabbit) *Container {

	userProfileRepo := repositories.NewUserProfileRepository(db)
	userProfileService := services.NewUserProfileService(userProfileRepo)
	userProfileController := controllers.NewUserProfileController(userProfileService, rabbit.Channel)

	followRepo := repositories.NewFollowRepository(db)
	followService := services.NewFollowService(followRepo)
	followController := controllers.NewFollowController(followService)

	return &Container{
		DB:     db,
		Rabbit: rabbit,
		// Controllers
		UserProfileController: userProfileController,
		FollowController:      followController,
		// Services
		UserProfileService: userProfileService,
		FollowService:      followService,
	}
}
