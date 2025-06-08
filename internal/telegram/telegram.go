package telegram

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/webshining/internal/common/database"
	"github.com/webshining/internal/common/rabbit"
	"github.com/webshining/internal/telegram/app"
	hndls "github.com/webshining/internal/telegram/handlers"
	"github.com/webshining/internal/telegram/middlewares"
	"github.com/webshining/internal/telegram/notifier"
)

type bot struct {
	bot        *gotgbot.Bot
	dispatcher *ext.Dispatcher
	db         *gorm.DB
	rabbit     *amqp.Connection
	logger     *zap.Logger
	notifier   *notifier.Notifier
}

func New() (*bot, error) {
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
	b, err := gotgbot.NewBot(os.Getenv("TELEGRAM_BOT_TOKEN"), nil)
	if err != nil {
		logger.Error("failed to create new bot:", zap.Error(err))
		return nil, err
	}
	dispatcher := ext.NewDispatcher(nil)

	// set bot commands
	commands := []gotgbot.BotCommand{
		{Command: "start", Description: "Start the bot"},
		{Command: "notify", Description: "Set channel notifications"},
	}
	if _, err := b.SetMyCommands(commands, nil); err != nil {
		logger.Error("failed to set bot commands", zap.Error(err))
		return nil, err
	}

	// setup app context
	app := app.New(b, rabbit, db, logger)

	// setup bot handlers
	middlewares := middlewares.New(app)
	hndl := hndls.New(app)

	dispatcher.AddHandlerToGroup(handlers.NewMessage(message.Command, middlewares.UserMiddleware), -10)
	dispatcher.AddHandlerToGroup(handlers.NewCallback(callbackquery.All, middlewares.UserMiddleware), -10)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("start", hndl.StartHandler), 10)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("notify", hndl.NotifyHandler), 10)
	dispatcher.AddHandlerToGroup(handlers.NewMessage(message.Audio, hndl.FileHandler), 10)
	dispatcher.AddHandlerToGroup(handlers.NewMessage(message.Voice, hndl.FileHandler), 10)
	dispatcher.AddHandlerToGroup(handlers.NewCallback(callbackquery.Prefix("guild:"), hndl.NotifyGuildHandler), 10)
	dispatcher.AddHandlerToGroup(handlers.NewCallback(callbackquery.Prefix("channel:"), hndl.NotifyChannelHandler), 10)

	// setup notifier
	notifier, err := notifier.New(app)
	if err != nil {
		logger.Error("failed to create notifier", zap.Error(err))
		return nil, err
	}

	return &bot{
		bot:        b,
		dispatcher: dispatcher,
		db:         db,
		rabbit:     rabbit,
		logger:     logger,
		notifier:   notifier,
	}, nil
}

func (b *bot) Run() {
	updater := ext.NewUpdater(b.dispatcher, nil)
	if err := updater.StartPolling(b.bot, nil); err != nil {
		b.logger.Error("failed to start polling", zap.Error(err))
		return
	}

	go func() {
		if err := b.notifier.Start(); err != nil {
			b.logger.Error("failed to start notifier", zap.Error(err))
		}
	}()

	b.logger.Info("Bot is now running. Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc
}
