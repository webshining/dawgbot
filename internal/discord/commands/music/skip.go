package music

import (
	"github.com/bwmarrin/discordgo"
)

func (m *Music) Skip(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Скипнул свагу",
		},
	})

	cancel, ok := m.playbackCancel[i.GuildID]
	if !ok {
		return
	}
	cancel()
	delete(m.playbackCancel, i.GuildID)
}
