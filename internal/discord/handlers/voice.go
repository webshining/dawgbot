package handlers

import (
	"encoding/json"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/rabbitmq/amqp091-go"
)

type VoiceJoinMessage struct {
	Username string `json:"username"`
	Channel  string `json:"channel"`
	Guild    string `json:"guild"`
}

func (h *Handlers) VoiceJoinHandler(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	if vs.ChannelID == "" {
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
		Username: user.Username,
		Channel:  channel.Name,
		Guild:    guild.Name,
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Printf("Ошибка при сериализации сообщения в JSON: %v\n", err)
		return
	}

	err = h.AMQP.Publish("", "voice", false, false,
		amqp091.Publishing{
			ContentType: "allpication/json",
			Body:        messageJSON,
		},
	)
	if err != nil {
		log.Printf("Ошибка при отправке сообщения в RabbitMQ: %v\n", err)
		return
	}
}
