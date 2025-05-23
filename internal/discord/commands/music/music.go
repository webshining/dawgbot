package music

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type Music struct {
	session        *discordgo.Session
	logger         *zap.Logger
	playbackCancel map[string]context.CancelFunc
}

func New(session *discordgo.Session, logger *zap.Logger) *Music {
	return &Music{
		session:        session,
		logger:         logger,
		playbackCancel: make(map[string]context.CancelFunc),
	}
}

func (m *Music) Commands() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	commands := make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
	commands["play"] = m.Play
	commands["skip"] = m.Skip
	return commands
}
