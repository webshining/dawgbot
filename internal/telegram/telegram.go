package telegram

import (
	"os"
	"os/signal"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"bot/internal/broker"
	"bot/internal/config"
	"bot/internal/database"
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
	logger     *zap.Logger
	notifier   *notifier.Notifier
}

func New() *bot {
	logger, _ := zap.NewDevelopment()
	config := config.MustLoad(logger)
	db := database.MustConnect(config.Database, logger)
	broker := broker.MustConnect("dawg-telegram", config.Broker, logger)

	// setup new bot session
	b, err := gotgbot.NewBot(config.Telegram.Token, nil)
	if err != nil {
		logger.Fatal("failed to create new bot:", zap.Error(err))
	}
	dispatcher := ext.NewDispatcher(nil)

	// setup bot commands
	commands := []gotgbot.BotCommand{
		{Command: "start", Description: "Start the bot"},
		{Command: "notify", Description: "Set channel notifications"},
	}
	if _, err := b.SetMyCommands(commands, nil); err != nil {
		logger.Fatal("failed to set bot commands", zap.Error(err))
	}

	// global context
	app := app.New(b, db, broker, logger)

	// modules
	start := start.New(app)
	notify := notify.New(app)
	user := user.New(app)

	// register modules
	registerHandler(dispatcher, 10, -10, user)
	registerHandler(dispatcher, 10, 0, start)
	registerHandler(dispatcher, 10, 0, notify)

	// setup notifier
	notifier := notifier.New(app)

	return &bot{
		bot:        b,
		dispatcher: dispatcher,
		db:         db,
		logger:     logger,
		notifier:   notifier,
	}
}

func (b *bot) Run() {
	updater := ext.NewUpdater(b.dispatcher, nil)
	if err := updater.StartPolling(b.bot, nil); err != nil {
		b.logger.Fatal("failed to start polling", zap.Error(err))
	}

	b.notifier.Start()

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
