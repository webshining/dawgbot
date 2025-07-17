package telegram

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"bot/internal/common/database"
	"bot/internal/common/rabbit"
	"bot/internal/telegram/app"
	"bot/internal/telegram/notifier"
	"bot/internal/telegram/notify"
	"bot/internal/telegram/start"
	"bot/internal/telegram/user"
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
	godotenv.Load()
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

	// setup bot commands
	commands := []gotgbot.BotCommand{
		{Command: "start", Description: "Start the bot"},
		{Command: "notify", Description: "Set channel notifications"},
	}
	if _, err := b.SetMyCommands(commands, nil); err != nil {
		logger.Error("failed to set bot commands", zap.Error(err))
		return nil, err
	}

	// global context
	app := app.New(b, rabbit, db, logger)

	// modules
	start := start.New(app)
	notify := notify.New(app)
	user := user.New(app)

	// register modules
	registerHandler(dispatcher, 10, -10, user)
	registerHandler(dispatcher, 10, 0, start)
	registerHandler(dispatcher, 10, 0, notify)

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

type handlerModule interface {
	Handlers(dp *ext.Dispatcher, group int)
	Middlewares(dp *ext.Dispatcher, group int)
}

func registerHandler(dp *ext.Dispatcher, group int, middlewaresGroup int, module handlerModule) {
	module.Middlewares(dp, middlewaresGroup)
	module.Handlers(dp, group)
}
