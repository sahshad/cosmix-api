package app

import (
	"log"
	"cosmix/shared/core/rabbitmq"
	consumer "notification-service/internal/messaging/consumers"
)

func RegisterConsumers(container *Container) {

	if err := rabbitmq.DeclareExchanges(container.Rabbit.Channel); err != nil {
		log.Fatal(err)
	}

	consumer.RegisterAuthUserRegisteredConsumer(
		container.Rabbit.Channel,
		container.EventService,
		)
}
