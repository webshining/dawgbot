package discord

import (
	"bot/internal/config"
	"bot/internal/discord/app"
	"bot/internal/discord/commands"
	"bot/internal/discord/handlers"
	"context"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type bot struct {
	session  *discordgo.Session
	logger   *zap.Logger
	commands []*discordgo.ApplicationCommand
}

func New(config *config.DiscordConfig, database *gorm.DB, telegramBot app.TelegramBot, logger *zap.Logger) *bot {
	// setup new bot session
	b, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		logger.Fatal("error creating bot session", zap.Error(err))
	}

	// setup app context
	app := app.New(b, database, telegramBot, logger)

	// set bot properties
	b.Identify.Intents = discordgo.IntentsGuildVoiceStates | discordgo.IntentsGuilds

	// register commands
	commands := commands.New(app)
	b.AddHandler(commands.Handler)

	// register handlers
	handlers := handlers.New(app, commands.Commands)
	for _, handler := range handlers.Handlers() {
		b.AddHandler(handler)
	}

	return &bot{
		session:  app.Session,
		logger:   app.Logger,
		commands: commands.Commands,
	}
}

func (b *bot) Run(ctx context.Context) {
	if err := b.session.Open(); err != nil {
		b.logger.Error("error opening connection to Discord", zap.Error(err))
		return
	}

	b.session.ApplicationCommandBulkOverwrite(b.session.State.User.ID, "", b.commands)

	b.logger.Info("Bot is now running")
	defer b.logger.Info("Bot is shuted down")

	<-ctx.Done()
}
