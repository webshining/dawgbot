package notifier

import (
	"encoding/json"
	"fmt"
	"html"
	"log"

	"github.com/PaulSonOfLars/gotgbot/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/webshining/internal/common/database"
	"github.com/webshining/internal/telegram/app"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Notifier struct {
	rabbit *amqp.Channel
	bot    *gotgbot.Bot
	db     *gorm.DB
	logger *zap.Logger
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
			var msg VoiceJoinMessage
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				log.Printf("error:%s", err)
				continue
			}

			var channel *database.Channel
			n.db.Preload("Users").First(&channel, "id = ?", msg.Channel)
			for _, user := range channel.Users {
				text := fmt.Sprintf("<code>[</code> <b>%s</b> <code>]</code> — <code>[</code> <b>%s</b> <code>]</code> — <code>[</code> <b>%s</b> <code>]</code>",
					html.EscapeString(msg.GuildName),
					html.EscapeString(msg.ChannelName),
					html.EscapeString(msg.Username),
				)
				if user.LastGuildID != msg.Guild {
					user.LastGuildID = msg.Guild
					n.db.Save(&user)
					n.bot.SendPhoto(user.TelegramID, gotgbot.InputFileByURL(msg.Image), &gotgbot.SendPhotoOpts{
						Caption:   text,
						ParseMode: "HTML",
					})
				} else {
					n.bot.SendMessage(user.TelegramID, text, &gotgbot.SendMessageOpts{
						ParseMode: "HTML",
					})
				}
			}
		}
	}()

	n.logger.Info("Notifier started")
	return nil
}
