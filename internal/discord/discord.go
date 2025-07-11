package discord

import (
	"fmt"
	"os"
	"os/signal"

	"bot/internal/common/database"
	"bot/internal/common/rabbit"
	"bot/internal/discord/app"
	"bot/internal/discord/commands"
	"bot/internal/discord/handlers"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Bot struct {
	session  *discordgo.Session
	rabbit   *amqp.Connection
	logger   *zap.Logger
	commands []*discordgo.ApplicationCommand
}

func New() (*Bot, error) {
	// load .env file
	godotenv.Load()

	// setup new logger
	logger, _ := zap.NewDevelopment()

	// setup new database connection
	db, err := database.New(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT")))
	if err != nil {
		logger.Error("error connecting to database", zap.Error(err))
		return nil, err
	}

	// setup new rabbit connection
	rabbit, err := rabbit.New(fmt.Sprintf("amqp://%s:%s@%s:%s/", os.Getenv("RB_USER"), os.Getenv("RB_PASS"), os.Getenv("RB_HOST"), os.Getenv("RB_PORT")))
	if err != nil {
		logger.Error("error connecting to rabbit", zap.Error(err))
		return nil, err
	}

	// setup new bot session
	bot, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		logger.Error("error creating bot session", zap.Error(err))
		return nil, err
	}

	// setup app context
	app := app.New(bot, db, rabbit, logger)

	// set bot properties
	bot.Identify.Intents = discordgo.IntentsGuildVoiceStates | discordgo.IntentsGuilds

	// register commands
	commands := commands.New(app)
	bot.AddHandler(commands.Handler)

	// register handlers
	handlers, err := handlers.New(app, commands.Commands)
	if err != nil {
		logger.Error("error creating handlers", zap.Error(err))
		return nil, err
	}
	for _, handler := range handlers.Handlers() {
		bot.AddHandler(handler)
	}

	return &Bot{
		session:  app.Session,
		rabbit:   app.Rabbit,
		logger:   app.Logger,
		commands: commands.Commands,
	}, nil
}

func (b *Bot) Run() {
	defer b.session.Close()
	defer b.rabbit.Close()
	defer b.logger.Sync()

	if err := b.session.Open(); err != nil {
		b.logger.Error("error opening connection to Discord", zap.Error(err))
		return
	}

	b.session.ApplicationCommandBulkOverwrite(b.session.State.User.ID, "", b.commands)

	b.logger.Info("Bot is now running. Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc
}
