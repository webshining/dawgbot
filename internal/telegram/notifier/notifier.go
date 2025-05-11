package notifier

import (
	"encoding/json"
	"fmt"
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
				n.Bot.SendMessage(user.ID, fmt.Sprintf("Пользователь **%s** присоединился к каналу **%s** на сервере **%s**", msg.Username, msg.ChannelName, msg.GuildName), &gotgbot.SendMessageOpts{ParseMode: "MarkdownV2"})
			}
		}
	}()

	n.Logger.Info("Notifier started")
	return nil
}
