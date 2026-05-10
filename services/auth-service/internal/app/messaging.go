package app

import (
	"cosmix/shared/core/rabbitmq"
	"log"
)

func RegisterConsumers(container *Container) {

	if err := rabbitmq.DeclareExchanges(container.Rabbit.Channel); err != nil {
		log.Fatal(err)
	}

}
