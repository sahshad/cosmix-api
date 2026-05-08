package messaging

import (
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Rabbit struct {
	URL        string
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func Connect(url string) *Rabbit {

	var conn *amqp.Connection
	var err error

	for {
		conn, err = amqp.Dial(url)
		if err == nil {
			break
		}

		log.Println("Waiting for RabbitMQ...")
		time.Sleep(5 * time.Second)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open channel:", err)
	}

	return &Rabbit{
		URL:        url,
		Connection: conn,
		Channel:    ch,
	}
}
