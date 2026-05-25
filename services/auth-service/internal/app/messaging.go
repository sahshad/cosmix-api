package app

import (
	"cosmix/shared/core/rabbitmq"
	// "log"
)

func RegisterSubscriptions(container *Container) error {

	if err := rabbitmq.DeclareExchanges(container.Rabbit.Channel); err != nil {
		return err
	}

	return nil
}
