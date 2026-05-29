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
		rabbitmq.PostAuthUserEmailVerificationCompletedQueue,
		rabbitmq.AuthUserEmailVerificationCompleted,
		container.PostUserSvc.HandleAuthUserEmailVerificationCompleted,
	); err != nil {
		return err
	}

	return nil
}
