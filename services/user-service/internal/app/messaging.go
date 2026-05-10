package app

import (
	"cosmix/shared/core/rabbitmq"
	"user-service/internal/messaging/consumer"
	"log"
)

func RegisterConsumers(container *Container) {

	if err := rabbitmq.DeclareExchanges(container.Rabbit.Channel); err != nil {
		log.Fatal(err)
	}
	
	consumer.RegisterAuthUserRegisteredConsumer(
		container.Rabbit.Channel,
		container.UserProfileService,
	)
}