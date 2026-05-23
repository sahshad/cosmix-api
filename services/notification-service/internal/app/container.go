package app

import (
	"cosmix/shared/core/rabbitmq"
	"notification-service/internal/controllers"
	grpcServer "notification-service/internal/grpc/server"
	"notification-service/internal/repositories"
	"notification-service/internal/services"

	"gorm.io/gorm"
)

type Container struct {
	DB     *gorm.DB
	Rabbit *rabbitmq.Rabbit

	// Controllers
	NotificationController *controllers.NotificationController

	// Services
	NotificationService           services.NotificationServiceInterface
	NotificationPreferenceService services.NotificationPreferenceServiceInterface
	EmailLogService               services.EmailLogServiceInterface
	EventService                  services.EventServiceInterface

	NotificationGrpcServer *grpcServer.NotificationServer
}

func NewContainer(db *gorm.DB, rabbit *rabbitmq.Rabbit) *Container {

	notificationRepo := repositories.NewNotificationRepository(db)
	notificationReferenceRepo := repositories.NewNotificationPreferenceRepository(db)
	emailLogRepo := repositories.NewEmailLogRepository(db)

	notificationSvc := services.NewNotificationService(notificationRepo)
	notificationReferenceSvc := services.NewNotificationPreferenceService(notificationReferenceRepo)
	emailLogSvc := services.NewEmailLogService(emailLogRepo)
	eventSvc := services.NewEventService(emailLogSvc, notificationSvc, notificationReferenceSvc)

	notificationCtrl := controllers.NewNotificationController(notificationSvc)

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
		NotificationController:        notificationCtrl,

		NotificationGrpcServer: notificationGrpcServer,
	}
}
