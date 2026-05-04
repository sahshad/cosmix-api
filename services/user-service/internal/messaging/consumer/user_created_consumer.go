package consumer

import (
	"encoding/json"
	"log"

	// "user-service/internal/dto"
	"user-service/internal/services"

	authEvents "cosmix-events/auth"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConsumeUserCreated(ch *amqp.Channel, userProfileService services.UserProfileServiceInterface) {

	q, err := ch.QueueDeclare(
		"user.created",
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
		"user.created",
		"auth.events",
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
			var event authEvents.UserCreated
			json.Unmarshal(msg.Body, &event)

			userCreatedEvent := authEvents.UserCreated{
				AuthUserID: event.AuthUserID,
				Email:      event.Email,
				FirstName:  event.FirstName,
				LastName:   event.LastName,
				CreatedAt:  event.CreatedAt,
			}

			err := userProfileService.CreateFromAuthEvent(userCreatedEvent)

			if err != nil {
				log.Println("User creation failed:", err)
			}
		}
	}()
}
