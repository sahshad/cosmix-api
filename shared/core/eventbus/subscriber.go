package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"go.uber.org/zap"
)

type EventHandler[T any] func(ctx context.Context, event T) error

func Subscribe[T any](
	ch *amqp.Channel,
	exchange string,
	queueName string,
	routingKey string,
	handler EventHandler[T],
) error {

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	err = ch.QueueBind(
		q.Name,
		routingKey,
		exchange,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	go func() {

		for msg := range msgs {

			start := time.Now()

			var event T

			if err := json.Unmarshal(
				msg.Body,
				&event,
			); err != nil {

				getLogger().Error(
					"Event Unmarshal Failed",
					zap.String(
						"routing_key",
						routingKey,
					),
					zap.String("error", err.Error()),
				)

				continue
			}

			ctx := context.Background()

			if value, ok := msg.Headers["correlation_id"]; ok {

				ctx = WithCorrelationID(
					ctx,
					fmt.Sprint(value),
				)
			}

			getLogger().Info(
				"Event Received",
				zap.String(
					"routing_key",
					routingKey,
				),
				zap.String(
					"correlation_id",
					CorrelationID(ctx),
				),
			)

			err := handler(
				ctx,
				event,
			)

			if err != nil {

				getLogger().Error(
					"Event Processing Failed",
					zap.String(
						"routing_key",
						routingKey,
					),
					zap.String(
						"correlation_id",
						CorrelationID(ctx),
					),
					zap.String("error", err.Error()),
				)

				continue
			}

			getLogger().Info(
				"Event Processed",
				zap.String(
					"routing_key",
					routingKey,
				),
				zap.String(
					"correlation_id",
					CorrelationID(ctx),
				),
				zap.Duration(
					"duration",
					time.Since(start),
				),
			)
		}
	}()

	return nil
}
