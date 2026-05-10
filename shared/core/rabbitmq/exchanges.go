package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	ExchangeEvents = "cosmix.events"
	ExchangeDLX    = "cosmix.events.dlx"
	ExchangeType   = "topic"
)

func DeclareExchanges(ch *amqp.Channel) error {

	// Exchange
	if err := ch.ExchangeDeclare(
		ExchangeEvents,
		ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	// DLX
	if err := ch.ExchangeDeclare(
		ExchangeDLX,
		ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	return nil
}
