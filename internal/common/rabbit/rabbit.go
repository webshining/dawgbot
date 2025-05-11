package rabbit

import amqp "github.com/rabbitmq/amqp091-go"

func New() (*amqp.Connection, *amqp.Channel, error) {
	amqp_conn, err := amqp.Dial("amqp://admin_user:admin_pass@localhost:5672/")
	if err != nil {
		return nil, nil, err
	}

	amqp_channel, err := amqp_conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	if _, err = amqp_channel.QueueDeclare("voice", true, false, false, false, nil); err != nil {
		return nil, nil, err
	}

	return amqp_conn, amqp_channel, nil
}
