package notifier

import (
	"bot/internal/common/database"
	"bot/internal/telegram/app"
	"encoding/json"
	"fmt"
	"html"

	"github.com/PaulSonOfLars/gotgbot/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Notifier struct {
	rabbit *amqp.Channel
	bot    *gotgbot.Bot
	db     *gorm.DB
	logger *zap.Logger
}

type Message struct {
	Message string `json:"message"`
	Data    []byte `json:"data"`
}

type VoiceJoinMessage struct {
	Username    string `json:"username"`
	Channel     string `json:"channel"`
	ChannelName string `json:"channel_name"`
	Guild       string `json:"guild"`
	GuildName   string `json:"guild_name"`
	Image       string `json:"image"`
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
		bot:    app.Bot,
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
			var msg Message
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				n.logger.Error("Failed to unmarshal message", zap.Error(err))
				continue
			}
			if msg.Message == "voice_join" {
				var data VoiceJoinMessage
				json.Unmarshal(msg.Data, &data)

				var channel *database.Channel
				n.db.Preload("Users").First(&channel, data.Channel)
				for _, user := range channel.Users {
					text := fmt.Sprintf("<code>[</code> <b>%s</b> <code>]</code> — <code>[</code> <b>%s</b> <code>]</code> — <code>[</code> <b>%s</b> <code>]</code>",
						html.EscapeString(data.GuildName),
						html.EscapeString(data.ChannelName),
						html.EscapeString(data.Username),
					)
					if user.LastGuildID != data.Guild {
						user.LastGuildID = data.Guild
						n.db.Save(&user)
						n.bot.SendPhoto(user.ID, gotgbot.InputFileByURL(data.Image), &gotgbot.SendPhotoOpts{
							Caption:   text,
							ParseMode: "HTML",
						})
					} else {
						n.bot.SendMessage(user.ID, text, &gotgbot.SendMessageOpts{
							ParseMode: "HTML",
						})
					}
				}
			}
		}
	}()

	n.logger.Info("Notifier started")
	return nil
}
