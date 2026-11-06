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
		rabbitmq.CommunicationAuthUserEmailVerificationCompletedQueue,
		rabbitmq.AuthUserEmailVerificationCompleted,
		container.CommunicationUserSvc.HandleAuthUserEmailVerificationCompleted,
	); err != nil {
		return err
	}

	if err := eventbus.Subscribe(
		container.Rabbit.Channel,
		rabbitmq.ExchangeEvents,
		rabbitmq.CommunicationUserProfileUpdatedQueue,
		rabbitmq.UserProfileUpdated,
		container.CommunicationUserSvc.HandleUserUpdated,
	); err != nil {
		return err
	}

	return nil
}
