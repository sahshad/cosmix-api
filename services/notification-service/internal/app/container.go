package app

import (
	"cosmix/shared/core/rabbitmq"
	// "notification-service/internal/contro.llers"
	grpcServer "notification-service/internal/grpc/server"
	"notification-service/internal/repositories"
	"notification-service/internal/services"

	"gorm.io/gorm"
)

type Container struct {
	DB     *gorm.DB
	Rabbit *rabbitmq.Rabbit

	// Controllers
	// NotificationController *controllers.NotificationController

	// Services
	NotificationService           *services.NotificationService
	NotificationPreferenceService *services.NotificationPreferenceService
	EmailLogService               *services.EmailLogService
	EventService                  *services.EventService
	NotificationUserService       *services.NotificationUserService

	NotificationGrpcServer *grpcServer.NotificationServer
}

func NewContainer(db *gorm.DB, rabbit *rabbitmq.Rabbit) *Container {

	notificationRepo := repositories.NewNotificationRepository(db)
	notificationReferenceRepo := repositories.NewNotificationPreferenceRepository(db)
	emailLogRepo := repositories.NewEmailLogRepository(db)
	notificationUserRepo := repositories.NewNotificationUserRepository(db)

	notificationSvc := services.NewNotificationService(notificationRepo)
	notificationReferenceSvc := services.NewNotificationPreferenceService(notificationReferenceRepo)
	emailLogSvc := services.NewEmailLogService(emailLogRepo)
	notificationUserSvc := services.NewNotificationUserService(notificationUserRepo)

	eventSvc := services.NewEventService(emailLogSvc, notificationSvc, notificationReferenceSvc, notificationUserSvc)

	// notificationCtrl := controllers.NewNotificationController(notificationSvc)

	notificationGrpcServer :=
		grpcServer.NewNotificationServer(
			notificationSvc,
		)

	return &Container{
		DB:                            db,
		Rabbit:                        rabbit,
		NotificationService:           notificationSvc,
		NotificationPreferenceService: notificationReferenceSvc,
		EmailLogService:               emailLogSvc,
		EventService:                  eventSvc,
		NotificationUserService:       notificationUserSvc,
		// NotificationController:        notificationCtrl,

		NotificationGrpcServer: notificationGrpcServer,
	}
}
