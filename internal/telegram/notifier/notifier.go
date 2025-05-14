package notifier

import (
	"encoding/json"
	"fmt"
	"html"
	"log"

	"github.com/PaulSonOfLars/gotgbot/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/webshining/internal/common/database"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Notifier struct {
	AMQP   *amqp.Channel
	Bot    *gotgbot.Bot
	DB     *gorm.DB
	Logger *zap.Logger
}

type VoiceJoinMessage struct {
	Username    string `json:"username"`
	Channel     string `json:"channel"`
	ChannelName string `json:"channel_name"`
	Guild       string `json:"guild"`
	GuildName   string `json:"guild_name"`
	Image       string `json:"image"`
}

func New(amqp *amqp.Channel, bot *gotgbot.Bot, db *gorm.DB, logger *zap.Logger) *Notifier {
	return &Notifier{
		AMQP:   amqp,
		Bot:    bot,
		DB:     db,
		Logger: logger,
	}
}

func (n *Notifier) Start() error {
	msgs, err := n.AMQP.Consume(
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
			n.DB.Preload("Users").First(&channel, "id = ?", msg.Channel)
			for _, user := range channel.Users {
				text := fmt.Sprintf("<code>[</code> <b>%s</b> <code>]</code> — <code>[</code> <b>%s</b> <code>]</code> — <code>[</code> <b>%s</b> <code>]</code>",
					html.EscapeString(msg.GuildName),
					html.EscapeString(msg.ChannelName),
					html.EscapeString(msg.Username),
				)
				if user.LastGuildID != msg.Guild {
					user.LastGuildID = msg.Guild
					n.DB.Save(&user)
					n.Bot.SendPhoto(user.TelegramID, gotgbot.InputFileByURL(msg.Image), &gotgbot.SendPhotoOpts{
						Caption:   text,
						ParseMode: "HTML",
					})
				} else {
					n.Bot.SendMessage(user.TelegramID, text, &gotgbot.SendMessageOpts{
						ParseMode: "HTML",
					})
				}
			}
		}
	}()

	n.Logger.Info("Notifier started")
	return nil
}
