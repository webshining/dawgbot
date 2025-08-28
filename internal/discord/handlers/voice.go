package handlers

import (
	"encoding/json"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type Message struct {
	Message string `json:"message"`
	Data    []byte `json:"data"`
}

type VoiceJoinMessage struct {
	Username string `json:"username"`
	Channel  string `json:"channel"`
	Guild    string `json:"guild"`
	Image    string `json:"image"`
}

func (h *handlers) VoiceJoinHandler(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	if vs.ChannelID == "" {
		return
	}
	if vs.BeforeUpdate != nil && vs.BeforeUpdate.ChannelID == vs.ChannelID {
		return
	}

	user, err := s.User(vs.UserID)
	if err != nil {
		return
	}

	channel, err := s.State.Channel(vs.ChannelID)
	if err != nil {
		channel, err = s.Channel(vs.ChannelID)
		if err != nil {
			return
		}
	}
	guild, err := s.State.Guild(vs.GuildID)
	if err != nil {
		guild, err = s.Guild(vs.GuildID)
		if err != nil {
			return
		}
	}

	message, err := json.Marshal(VoiceJoinMessage{
		Username: user.DisplayName(),
		Channel:  channel.ID,
		Guild:    guild.ID,
		Image:    guild.IconURL("1024"),
	})
	if err != nil {
		h.app.Logger.Error("failed to marshal message", zap.Error(err))
		return
	}

	if token := h.app.Broker.Publish("voice", message); token.Wait() && token.Error() != nil {
		h.app.Logger.Error("failed to publish message", zap.Error(token.Error()))
	}
}
