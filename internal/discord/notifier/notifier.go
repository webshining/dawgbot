package notifier

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/webshining/internal/discord/app"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Notifier struct {
	rabbit *amqp.Channel
	db     *gorm.DB
	logger *zap.Logger
}

func New(app *app.AppContext) (*Notifier, error) {
	rabbitChannel, err := app.Rabbit.Channel()
	if err != nil {
		return nil, err
	}
	if _, err = rabbitChannel.QueueDeclare("voice", true, false, false, false, nil); err != nil {
		return nil, err
	}
	return &Notifier{
		rabbit: rabbitChannel,
		db:     app.DB,
		logger: app.Logger,
	}, nil
}

func (n *Notifier) Start() error {
	msgs, err := n.rabbit.Consume(
		"voice",
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
		for d := range msgs {
			var msg any
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				n.logger.Error("Failed to unmarshal message", zap.Error(err))
				continue
			}
		}
	}()

	n.logger.Info("Notifier started")
	return nil
}
