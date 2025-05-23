package rabbit

import amqp "github.com/rabbitmq/amqp091-go"

func New(url string) (*amqp.Channel, error) {
	amqp_conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	amqp_channel, err := amqp_conn.Channel()
	if err != nil {
		return nil, err
	}

	if _, err = amqp_channel.QueueDeclare("voice", true, false, false, false, nil); err != nil {
		return nil, err
	}

	return amqp_channel, nil
}
