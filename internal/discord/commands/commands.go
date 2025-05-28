package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/webshining/internal/discord/app"
	"github.com/webshining/internal/discord/commands/music"
	"github.com/webshining/internal/discord/commands/notify"
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

	musicModule := music.New(app)
	notifyModule := notify.New(app)
	c.registerModules(musicModule, notifyModule)

	return c
}

func (c *commands) registerModules(modules ...commandModule) *commands {
	for _, module := range modules {
		c.Commands = append(c.Commands, module.Definitions()...)
		for name, handler := range module.Commands() {
			c.handlers[name] = handler
		}
	}
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
