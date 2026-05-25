package app

import (
	"cosmix/shared/core/eventbus"
	"cosmix/shared/core/rabbitmq"
	// "log"
	// "post-service/internal/messaging/consumer"
)

func RegisterSubscriptions(container *Container) error {

	if err := rabbitmq.DeclareExchanges(container.Rabbit.Channel); err != nil {
		return err
	}

	if err := eventbus.Subscribe(
		container.Rabbit.Channel,
		rabbitmq.ExchangeEvents,
		rabbitmq.PostAuthUserRegistered,
		rabbitmq.AuthUserRegistered,
		container.PostUserSvc.CreateFromAuthEvent,
	); err != nil {
		return err
	}

	// consumer.RegisterAuthUserRegisteredConsumer(
	// 	container.Rabbit.Channel,
	// 	container.PostUserSvc,
	// )
	return nil
}
