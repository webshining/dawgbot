package notify

import (
	"github.com/bwmarrin/discordgo"
	"github.com/webshining/internal/discord/app"
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
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Notifications have been enabled for this guild.",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label: "Telegram",
							Style: discordgo.LinkButton,
							URL:   "https://t.me/dawgdsbot?start=" + i.GuildID,
						},
					},
				},
			},
		},
	}); err != nil {
		n.logger.Error("Failed to respond to interaction", zap.Error(err), zap.String("guild_id", i.GuildID))
		return
	}
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
