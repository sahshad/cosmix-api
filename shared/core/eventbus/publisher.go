package eventbus

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

func Publish(
	ctx context.Context,
	ch *amqp.Channel,
	exchange string,
	routingKey string,
	payload any,
) error {

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	headers := amqp.Table{}

	if correlationID := CorrelationID(ctx); correlationID != "" {
		headers["correlation_id"] = correlationID
	}

	err = ch.PublishWithContext(
		ctx,
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Headers:     headers,
		},
	)

	if err != nil {

		getLogger().Error(
			"event publish failed",
			zap.String("routing_key", routingKey),
			zap.Error(err),
		)

		return fmt.Errorf(
			"publish %s: %w",
			routingKey,
			err,
		)
	}

	getLogger().Info(
		"event published",
		zap.String("routing_key", routingKey),
		zap.String(
			"correlation_id",
			CorrelationID(ctx),
		),
	)

	return nil
}
