package app

import (
	"auth-service/internal/messaging"
	"log"
)

func RegisterConsumers(container *Container) {

	if err := messaging.DeclareExchanges(container.Rabbit.Channel); err != nil {
		log.Fatal(err)
	}
	
}
