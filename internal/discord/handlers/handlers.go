package handlers

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/webshining/internal/discord/app"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type handlers struct {
	db     *gorm.DB
	rabbit *amqp.Channel
	logger *zap.Logger
}

func New(app *app.AppContext) (*handlers, error) {
	rabbitChannel, err := app.Rabbit.Channel()
	if err != nil {
		return nil, err
	}
	if _, err = rabbitChannel.QueueDeclare("voice", true, false, false, false, nil); err != nil {
		return nil, err
	}

	return &handlers{db: app.DB, rabbit: rabbitChannel, logger: app.Logger}, nil
}

func (h *handlers) Handlers() []interface{} {
	hdls := []interface{}{
		h.GuildAddHandler,
		h.GuildUpdateHandler,
		h.GuildDeleteHandler,

		h.VoiceJoinHandler,

		h.ChannelAddHandler,
		h.ChannelUpdateHandler,
		h.ChannelDeleteHandler,
	}
	return hdls
}
