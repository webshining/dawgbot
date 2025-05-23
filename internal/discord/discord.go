package discord

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/webshining/internal/common/database"
	"github.com/webshining/internal/common/rabbit"
	"github.com/webshining/internal/discord/commands"
	"github.com/webshining/internal/discord/handlers"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Bot struct {
	session  *discordgo.Session
	db       *gorm.DB
	rabbit   *amqp.Channel
	logger   *zap.Logger
	commands []*discordgo.ApplicationCommand
}

func New() *Bot {
	// load .env file
	godotenv.Load()

	// setup new logger
	logger, _ := zap.NewDevelopment()

	// setup new database connection
	db, err := database.New(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT")))
	if err != nil {
		logger.Error("error connecting to database", zap.Error(err))
		return nil
	}

	// setup new rabbit connection
	amqp_ch, err := rabbit.New(fmt.Sprintf("amqp://%s:%s@%s:%s/", os.Getenv("RB_USER"), os.Getenv("RB_PASS"), os.Getenv("RB_HOST"), os.Getenv("RB_PORT")))
	if err != nil {
		logger.Error("error connecting to rabbit", zap.Error(err))
		return nil
	}

	// setup new bot session
	bot, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		logger.Error("error creating bot session", zap.Error(err))
		return nil
	}

	// setup bot handlers
	bot.Identify.Intents = discordgo.IntentsGuildVoiceStates | discordgo.IntentsGuilds
	handlers := handlers.New(db, amqp_ch, logger)
	commands := commands.New(bot, logger)

	// voice handlers
	bot.AddHandler(handlers.VoiceJoinHandler)
	// guild handlers
	bot.AddHandler(handlers.GuildAddHandler)
	bot.AddHandler(handlers.GuildUpdateHandler)
	bot.AddHandler(handlers.GuildDeleteHandler)
	// channel handlers
	bot.AddHandler(handlers.ChannelAddHandler)
	bot.AddHandler(handlers.ChannelUpdateHandler)
	bot.AddHandler(handlers.ChannelDeleteHandler)
	// commands handlers
	bot.AddHandler(commands.Handler)

	return &Bot{
		session:  bot,
		db:       db,
		rabbit:   amqp_ch,
		logger:   logger,
		commands: commands.Commands,
	}
}

func (b *Bot) Run() {
	defer b.session.Close()
	defer b.rabbit.Close()
	defer b.logger.Sync()

	if err := b.session.Open(); err != nil {
		b.logger.Error("error opening connection to Discord", zap.Error(err))
		return
	}

	// register commands for all guilds
	guilds, _ := b.session.UserGuilds(100, "", "", false)
	for _, guild := range guilds {
		b.session.ApplicationCommandBulkOverwrite(b.session.State.User.ID, guild.ID, b.commands)
	}

	b.logger.Info("Bot is now running. Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc
}
