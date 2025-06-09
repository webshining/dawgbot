package music

import (
	"context"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/webshining/internal/discord/app"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Music struct {
	playbackCancel       map[string]context.CancelFunc
	autoDisconnectTimers map[string]*time.Timer

	session *discordgo.Session
	logger  *zap.Logger
	db      *gorm.DB
}

func New(app *app.AppContext) *Music {
	return &Music{
		session:              app.Session,
		logger:               app.Logger,
		db:                   app.DB,
		playbackCancel:       make(map[string]context.CancelFunc),
		autoDisconnectTimers: make(map[string]*time.Timer),
	}
}

func (m *Music) Commands() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	commands := make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
	commands["play"] = m.Play
	commands["skip"] = m.Skip
	return commands
}

func (m *Music) Definitions() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "play",
			Description: "Play a file",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "file",
					Description: "The file to play",
					Type:        discordgo.ApplicationCommandOptionAttachment,
				},
				{
					Name:        "youtubeurl",
					Description: "The youtube url to play",
					Type:        discordgo.ApplicationCommandOptionString,
				},
			},
		},
		{
			Name:        "skip",
			Description: "skip current song",
		},
	}
}
