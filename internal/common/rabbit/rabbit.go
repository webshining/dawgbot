package rabbit

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func New(url string) (*amqp.Connection, error) {
	var (
		amqpConn *amqp.Connection
		err      error
	)

	timeout := time.After(1 * time.Minute)
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		amqpConn, err = amqp.Dial(url)
		if err == nil {
			return amqpConn, nil
		}

		select {
		case <-timeout:
			return nil, err
		case <-ticker.C:
			// next try
		}
	}
}
