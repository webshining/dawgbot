package discord

import (
	"os"
	"os/signal"

	"bot/internal/broker"
	"bot/internal/config"
	"bot/internal/database"
	"bot/internal/discord/app"
	"bot/internal/discord/commands"
	"bot/internal/discord/handlers"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type Bot struct {
	session  *discordgo.Session
	logger   *zap.Logger
	commands []*discordgo.ApplicationCommand
}

func New() *Bot {
	// setup new logger
	logger, _ := zap.NewDevelopment()
	config := config.MustLoad(logger)
	db := database.MustConnect(config.Database, logger)
	broker := broker.MustConnect("dawg-discord", config.Broker, logger)

	// setup new bot session
	bot, err := discordgo.New("Bot " + config.Discord.Token)
	if err != nil {
		logger.Named("discord").Fatal("error creating bot session", zap.Error(err))
	}

	// setup app context
	app := app.New(bot, db, broker, logger)

	// set bot properties
	bot.Identify.Intents = discordgo.IntentsGuildVoiceStates | discordgo.IntentsGuilds

	// register commands
	commands := commands.New(app)
	bot.AddHandler(commands.Handler)

	// register handlers
	handlers := handlers.New(app, commands.Commands)
	for _, handler := range handlers.Handlers() {
		bot.AddHandler(handler)
	}

	return &Bot{
		session:  app.Session,
		logger:   app.Logger,
		commands: commands.Commands,
	}
}

func (b *Bot) Run() {
	defer b.session.Close()
	defer b.logger.Sync()

	if err := b.session.Open(); err != nil {
		b.logger.Fatal("error opening connection to Discord", zap.Error(err))
	}

	b.session.ApplicationCommandBulkOverwrite(b.session.State.User.ID, "", b.commands)

	b.logger.Info("Bot is now running. Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc
}
