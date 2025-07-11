package commands

import (
	"bot/internal/discord/app"
	"bot/internal/discord/commands/notify"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type commandModule interface {
	Commands() map[string]func(*discordgo.Session, *discordgo.InteractionCreate)
	Definitions() []*discordgo.ApplicationCommand
}

type commands struct {
	Commands []*discordgo.ApplicationCommand
	handlers map[string]func(*discordgo.Session, *discordgo.InteractionCreate)

	session *discordgo.Session
	logger  *zap.Logger
}

func New(app *app.AppContext) *commands {
	c := &commands{
		session:  app.Session,
		logger:   app.Logger,
		handlers: make(map[string]func(*discordgo.Session, *discordgo.InteractionCreate)),
	}

	notifyModule := notify.New(app)
	c.registerModules(notifyModule)

	return c
}

func (c *commands) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cmd, ok := c.handlers[i.ApplicationCommandData().Name]
	if !ok {
		c.logger.Error("command not found", zap.String("command", i.ApplicationCommandData().Name))
		return
	}

	cmd(s, i)
}
