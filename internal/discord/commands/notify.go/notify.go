package notify

import (
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type Notify struct {
	session *discordgo.Session
	logger  *zap.Logger
}

func New(session *discordgo.Session, logger *zap.Logger) *Notify {
	return &Notify{
		session: session,
		logger:  logger,
	}
}

func (n *Notify) Commands() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	commands := make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
	commands["notify"] = n.NotifyHandler
	return commands
}

func (n *Notify) NotifyHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
