package publisher

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	authEvents "cosmix/shared/events/auth"
	"cosmix/shared/core/rabbitmq"
)

func PublishAuthUserRegistered(ch *amqp.Channel, event authEvents.AuthUserRegistered) {
	body, _ := json.Marshal(event)

	err := ch.PublishWithContext(
		context.Background(),
		rabbitmq.ExchangeEvents,
		rabbitmq.AuthUserRegistered,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		log.Println("Failed to publish user.created event:", err)
	}
}
