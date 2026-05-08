package app

import (
	"log"
	"notification-service/internal/messaging"
	consumer "notification-service/internal/messaging/consumers"
)

func RegisterConsumers(container *Container) {

	if err := messaging.DeclareExchanges(container.Rabbit.Channel); err != nil {
		log.Fatal(err)
	}

	consumer.RegisterAuthUserRegisteredConsumer(
		container.Rabbit.Channel,
		container.EventService,
		)
}
