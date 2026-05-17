package consumer

import (
	"encoding/json"
	"log"

	"post-service/internal/services"

	"cosmix/shared/core/rabbitmq"
	authEvents "cosmix/shared/events/auth"

	amqp "github.com/rabbitmq/amqp091-go"
)

func RegisterAuthUserRegisteredConsumer(ch *amqp.Channel, postUserService services.PostUserServiceInterface) {

	q, err := ch.QueueDeclare(
		rabbitmq.PostUserAuthUserRegistered,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Queue declaration failed:", err)
	}

	if err := ch.QueueBind(
		q.Name,
		rabbitmq.AuthUserRegistered,
		rabbitmq.ExchangeEvents,
		false,
		nil,
	); err != nil {
		log.Fatal("Queue binding failed:", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true, // auto-ack (OK for MVP)
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal("Consumer failed:", err)
	}

	go func() {
		for msg := range msgs {
			var event authEvents.AuthUserRegistered
			json.Unmarshal(msg.Body, &event)

			userCreatedEvent := authEvents.AuthUserRegistered{
				AuthUserID:  event.AuthUserID,
				Email:       event.Email,
				Username:    event.Username,
				DisplayName: event.DisplayName,
				CreatedAt:   event.CreatedAt,
			}

			err := postUserService.CreateFromAuthEvent(userCreatedEvent)

			if err != nil {
				log.Println("User creation failed:", err)
			}
		}
	}()
}
