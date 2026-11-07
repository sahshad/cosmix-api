package app

import (
	"communication-service/internal/grpc"
	"communication-service/internal/repositories"
	"communication-service/internal/services"

	"cosmix/shared/core/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type Container struct {
	DB     *gorm.DB
	Rabbit *rabbitmq.Rabbit

	ConversationService  *services.ConversationService
	MessageService       *services.MessageService
	CommunicationUserSvc *services.CommunicationUserService

	CommunicationGrpcServer *grpc.CommunicationServer
}

func NewContainer(db *gorm.DB, rabbit *rabbitmq.Rabbit) *Container {
	var rabbitCh *amqp.Channel
	if rabbit != nil {
		rabbitCh = rabbit.Channel
	}

	communicationUserRepo := repositories.NewCommunicationUserRepository(db)
	conversationRepo := repositories.NewConversationRepository(db)
	messageRepo := repositories.NewMessageRepository(db)
	receiptRepo := repositories.NewMessageReadReceiptRepository(db)

	communicationUserService := services.NewCommunicationUserService(communicationUserRepo)
	conversationService := services.NewConversationService(conversationRepo, communicationUserRepo, receiptRepo)
	messageService := services.NewMessageService(messageRepo, conversationRepo, communicationUserRepo, receiptRepo, rabbitCh)

	communicationGrpcServer := grpc.NewCommunicationServer(
		conversationService,
		messageService,
	)

	return &Container{
		DB:     db,
		Rabbit: rabbit,

		ConversationService:  conversationService,
		MessageService:       messageService,
		CommunicationUserSvc: communicationUserService,

		CommunicationGrpcServer: communicationGrpcServer,
	}
}
