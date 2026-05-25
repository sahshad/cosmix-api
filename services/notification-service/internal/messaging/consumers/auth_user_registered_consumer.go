package consumer

// import (
// 	"encoding/json"
// 	"log"

// 	"notification-service/internal/services"

// 	authEvents "cosmix/shared/events/auth"
// 	"cosmix/shared/core/rabbitmq"

// 	amqp "github.com/rabbitmq/amqp091-go"
// )

// func RegisterAuthUserRegisteredConsumer(
// 	ch *amqp.Channel,
// 	eventSvc services.EventServiceInterface,
// ) {

// 	q, err := ch.QueueDeclare(
// 		rabbitmq.NotificationAuthUserRegistered,
// 		true,
// 		false,
// 		false,
// 		false,
// 		nil,
// 	)
// 	if err != nil {
// 		log.Fatal("Queue declaration failed:", err)
// 	}

// 	if err := ch.QueueBind(
// 		q.Name,
// 		rabbitmq.AuthUserRegistered,
// 		rabbitmq.ExchangeEvents,
// 		false,
// 		nil,
// 	); err != nil {
// 		log.Fatal("Queue binding failed:", err)
// 	}

// 	msgs, err := ch.Consume(
// 		q.Name,
// 		"",
// 		true, // auto-ack (OK for MVP)
// 		false,
// 		false,
// 		false,
// 		nil,
// 	)

// 	if err != nil {
// 		log.Fatal("Consumer failed:", err)
// 	}

// 	go func() {
// 		for msg := range msgs {
// 			var event authEvents.AuthUserRegistered
// 			if err := json.Unmarshal(msg.Body, &event); err != nil {
// 				log.Println("Auth user registerd event body parsing failed:", err)
// 			}

// 			log.Println("listening user created event from notification service")
// 			if err := eventSvc.HandleUserRegistered(event); err != nil {
// 				log.Println("User registerd event handling failed:", err)
// 			}
// 		}
// 	}()
// }
