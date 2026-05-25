package app

import (
	"cosmix/shared/core/rabbitmq"
	// "user-service/internal/messaging/consumer"
	"cosmix/shared/core/eventbus"
	// "log"
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
	// consumer.RegisterAuthUserRegisteredConsumer(
	// 	container.Rabbit.Channel,
	// 	container.UserProfileService,
	// )

	return nil
}