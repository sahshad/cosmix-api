package app

import (
	"cosmix/shared/core/eventbus"
	"cosmix/shared/core/rabbitmq"
)

func RegisterSubscriptions(container *Container) error {
	if err := rabbitmq.DeclareExchanges(container.Rabbit.Channel); err != nil {
		return err
	}

	if err := eventbus.Subscribe(
		container.Rabbit.Channel,
		rabbitmq.ExchangeEvents,
		rabbitmq.UserAuthUserRegistered,
		rabbitmq.AuthUserRegistered,
		container.UserService.CreateFromAuthEvent,
	); err != nil {
		return err
	}

	return nil
}
