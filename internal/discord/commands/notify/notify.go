package notify

import (
	"fmt"

	"bot/internal/discord/app"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type Notify struct {
	session *discordgo.Session
	logger  *zap.Logger
}

func New(app *app.AppContext) *Notify {
	return &Notify{
		session: app.Session,
		logger:  app.Logger,
	}
}

func (n *Notify) notifyHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "To receive notifications in Telegram, please click the button below.",
			Flags:   discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label: "Telegram",
							Style: discordgo.LinkButton,
							URL:   fmt.Sprintf("https://t.me/dawgdsbot?start=%s_%s", i.GuildID, i.Member.User.ID),
						},
					},
				},
			},
		},
	})
}

func (n *Notify) Commands() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	commands := make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
	commands["notify"] = n.notifyHandler
	return commands
}

func (m *Notify) Definitions() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "notify",
			Description: "allow notifications in telegram for this guild",
		},
	}
}
