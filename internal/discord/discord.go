package discord

import (
	"fmt"
	"os"

	"bot/internal/common/database"
	"bot/internal/discord/app"
	"bot/internal/discord/commands"
	"bot/internal/discord/handlers"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type bot struct {
	session  *discordgo.Session
	logger   *zap.Logger
	commands []*discordgo.ApplicationCommand
}

func New(telegramBot *gotgbot.Bot) *bot {
	// load .env file
	godotenv.Load()

	// setup new logger
	logger, _ := zap.NewDevelopment()

	// setup new database connection
	db, err := database.New(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT")))
	if err != nil {
		logger.Fatal("error connecting to database", zap.Error(err))
	}

	// setup new bot session
	b, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		logger.Fatal("error creating bot session", zap.Error(err))
	}

	// setup app context
	app := app.New(b, db, telegramBot, logger)

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

func (b *bot) Run() {
	if err := b.session.Open(); err != nil {
		b.logger.Error("error opening connection to Discord", zap.Error(err))
		return
	}

	b.session.ApplicationCommandBulkOverwrite(b.session.State.User.ID, "", b.commands)

	b.logger.Info("Bot is now running")
}
