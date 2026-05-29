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
		rabbitmq.NotificationAuthUserEmailVerificationCompletedQueue,
		rabbitmq.AuthUserEmailVerificationCompleted,
		container.EventService.HandleUserEmailVerificationCompleted,
	); err != nil {
		return err
	}

	if err := eventbus.Subscribe(
		container.Rabbit.Channel,
		rabbitmq.ExchangeEvents,
		rabbitmq.NotificationAuthUserEmailVerificationRequestedQueue,
		rabbitmq.AuthUserEmailVerificationRequested,
		container.NotificationService.HandleEmailVerification,
	); err != nil {
		return err
	}

	if err := eventbus.Subscribe(
		container.Rabbit.Channel,
		rabbitmq.ExchangeEvents,
		rabbitmq.NotificationAuthUserForgotPasswordRequestedQueue,
		rabbitmq.AuthUserForgotPasswordRequested,
		container.NotificationService.HandleAuthUserForgotPasswordRequest,
	); err != nil {
		return err
	}

	if err := eventbus.Subscribe(
		container.Rabbit.Channel,
		rabbitmq.ExchangeEvents,
		rabbitmq.NotificationAuthUserPasswordChangedQueue,
		rabbitmq.AuthUserPasswordChanged,
		container.NotificationService.HandleAuthUserPasswordChanged,
	); err != nil {
		return err
	}

	return nil
}
