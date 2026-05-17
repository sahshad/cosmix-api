package app

import (
	"cosmix/shared/core/rabbitmq"
	"log"
	"post-service/internal/messaging/consumer"
)

func RegisterConsumers(container *Container) {

	if err := rabbitmq.DeclareExchanges(container.Rabbit.Channel); err != nil {
		log.Fatal(err)
	}

	consumer.RegisterAuthUserRegisteredConsumer(
		container.Rabbit.Channel,
		container.PostUserSvc,
	)
}
