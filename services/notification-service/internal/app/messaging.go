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
		rabbitmq.NotificationAuthUserRegistered,
		rabbitmq.AuthUserRegistered,
		container.EventService.HandleUserRegistered,
	); err != nil {
		return err
	}

	if err := eventbus.Subscribe(
		container.Rabbit.Channel,
		rabbitmq.ExchangeEvents,
		rabbitmq.NotificationAuthUserEmailVerification,
		rabbitmq.AuthUserEmailVerification,
		container.NotificationService.HandleEmailVerification,
	); err != nil {
		return err
	}

	return nil
}
