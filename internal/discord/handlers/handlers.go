package handlers

import (
	"bot/internal/discord/app"

	"github.com/bwmarrin/discordgo"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type handlers struct {
	db       *gorm.DB
	rabbit   *amqp.Channel
	logger   *zap.Logger
	commands []*discordgo.ApplicationCommand
}

func New(app *app.AppContext, commands []*discordgo.ApplicationCommand) (*handlers, error) {
	rabbitChannel, err := app.Rabbit.Channel()
	if err != nil {
		return nil, err
	}
	if _, err = rabbitChannel.QueueDeclare("voice", true, false, false, false, nil); err != nil {
		return nil, err
	}

	return &handlers{db: app.DB, rabbit: rabbitChannel, logger: app.Logger, commands: commands}, nil
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
