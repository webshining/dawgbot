package handlers

import (
	"encoding/json"

	"github.com/bwmarrin/discordgo"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type VoiceJoinMessage struct {
	Username    string `json:"username"`
	Channel     string `json:"channel"`
	ChannelName string `json:"channel_name"`
	Guild       string `json:"guild"`
	GuildName   string `json:"guild_name"`
}

func (h *Handlers) VoiceJoinHandler(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
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

	message := VoiceJoinMessage{
		Username:    user.Username,
		Channel:     channel.ID,
		ChannelName: channel.Name,
		Guild:       guild.ID,
		GuildName:   guild.Name,
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		h.Logger.Error("failed to marshal message", zap.Error(err))
		return
	}

	err = h.AMQP.Publish("", "voice", false, false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        messageJSON,
		},
	)
	if err != nil {
		h.Logger.Error("failed to publish message", zap.Error(err))
	}
}
