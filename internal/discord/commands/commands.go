package commands

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/webshining/internal/discord/commands/music"
	"github.com/webshining/internal/discord/commands/notify.go"
	"go.uber.org/zap"
)

type Command struct {
	Name    string
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

type Commands struct {
	Session  *discordgo.Session
	Commands []*discordgo.ApplicationCommand

	logger         *zap.Logger
	playbackCancel map[string]context.CancelFunc
	handlers       map[string]func(*discordgo.Session, *discordgo.InteractionCreate)
}

func New(session *discordgo.Session, logger *zap.Logger) *Commands {
	commands := &Commands{
		Session: session,
		Commands: []*discordgo.ApplicationCommand{
			{
				Name:        "play",
				Description: "Play a file",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Name:        "file",
						Description: "The file to play",
						Type:        discordgo.ApplicationCommandOptionAttachment,
						Required:    true,
					},
				},
			},
			{
				Name:        "skip",
				Description: "skip current song",
			},
			{
				Name:        "notify",
				Description: "allow notifications in telegram for this guild",
			},
		},

		logger:         logger,
		playbackCancel: make(map[string]context.CancelFunc),
		handlers:       make(map[string]func(*discordgo.Session, *discordgo.InteractionCreate)),
	}

	// Initialize command handlers
	music := music.New(session, logger)
	notify := notify.New(session, logger)

	// Register commands from different modules
	commands.Register(notify.Commands())
	commands.Register(music.Commands())

	return commands
}

func (c *Commands) Register(commands map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	for name, cmd := range commands {
		c.handlers[name] = cmd
	}
}

func (c *Commands) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cmd, ok := c.handlers[i.ApplicationCommandData().Name]
	if !ok {
		c.logger.Error("command not found", zap.String("command", i.ApplicationCommandData().Name))
		return
	}

	cmd(s, i)
}
