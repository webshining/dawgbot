package notifier

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/PaulSonOfLars/gotgbot/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/webshining/internal/common/database"
	"gorm.io/gorm"
)

type Notifier struct {
	AMQP *amqp.Channel
	Bot  *gotgbot.Bot
	DB   *gorm.DB
}

type VoiceJoinMessage struct {
	Username string `json:"username"`
	Channel  string `json:"channel"`
	Guild    string `json:"guild"`
}

func New(amqp *amqp.Channel, bot *gotgbot.Bot, db *gorm.DB) *Notifier {
	return &Notifier{
		AMQP: amqp,
		Bot:  bot,
		DB:   db,
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

			var users []*database.User
			n.DB.Find(&users)
			for _, user := range users {
				n.Bot.SendMessage(user.ID, fmt.Sprintf("Пользователь %s присоединился к каналу %s на сервере %s", msg.Username, msg.Channel, msg.Guild), nil)
			}
		}
	}()

	log.Println(" [*] Консьюмер voice запущен и слушает очередь.")
	return nil
}
