package notifier

import (
	"bot/internal/common/database"
	"bot/internal/telegram/app"
	"encoding/json"
	"fmt"
	"html"

	"github.com/PaulSonOfLars/gotgbot/v2"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

type Notifier struct {
	app *app.AppContext
}

type VoiceJoinMessage struct {
	Username string `json:"username"`
	Channel  string `json:"channel"`
	Guild    string `json:"guild"`
	Image    string `json:"image"`
}

func New(app *app.AppContext) *Notifier {
	return &Notifier{app}
}

func (n *Notifier) handleMessage(client mqtt.Client, msg mqtt.Message) {
	var data VoiceJoinMessage
	json.Unmarshal(msg.Payload(), &data)

	var channel *database.Channel
	var guild *database.Guild
	n.app.DB.Preload("Users").First(&channel, data.Channel)
	n.app.DB.First(&guild, data.Guild)
	for _, user := range channel.Users {
		text := fmt.Sprintf("<code>[</code> <b>%s</b> <code>]</code> — <code>[</code> <b>%s</b> <code>]</code> — <code>[</code> <b>%s</b> <code>]</code>",
			html.EscapeString(guild.Name),
			html.EscapeString(channel.Name),
			html.EscapeString(data.Username),
		)
		if user.LastGuildID != guild.ID {
			user.LastGuildID = data.Guild
			n.app.DB.Save(&user)
			n.app.Bot.SendPhoto(user.ID, gotgbot.InputFileByURL(data.Image), &gotgbot.SendPhotoOpts{
				Caption:   text,
				ParseMode: "HTML",
			})
		} else {
			n.app.Bot.SendMessage(user.ID, text, &gotgbot.SendMessageOpts{
				ParseMode: "HTML",
			})
		}
	}
}

func (n *Notifier) Start() {
	if token := n.app.Broker.Subscribe("voice", n.handleMessage); token.Wait() && token.Error() != nil {
		n.app.Logger.Error("failed to recieve message", zap.Error(token.Error()))
	}

	n.app.Logger.Info("Notifier started")
}
