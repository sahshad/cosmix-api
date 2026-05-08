package messaging

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"cosmix-events/rabbitmq"
)


func DeclareExchanges(ch *amqp.Channel) error {

	// Exchange
	if err := ch.ExchangeDeclare(
		rabbitmq.ExchangeEvents,
		rabbitmq.ExchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	//  DLX
	if err := ch.ExchangeDeclare(
		rabbitmq.ExchangeDLX,
		rabbitmq.ExchangeType,
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
