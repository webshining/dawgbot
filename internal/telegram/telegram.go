package telegram

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	hndls "github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"bot/internal/config"
	"bot/internal/database"
	"bot/internal/telegram/handlers"
	"bot/internal/telegram/middlewares"
)

type bot struct {
	Bot        *gotgbot.Bot
	dispatcher *ext.Dispatcher
	logger     *zap.Logger
}

func New() *bot {
	godotenv.Load()
	logger, _ := zap.NewDevelopment()
	config := config.MustLoad(logger)

	// setup new database connection
	db := database.MustConnect(config.Database, logger)

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

	// middlewares
	middlewares := middlewares.New(db, logger)
	dispatcher.AddHandlerToGroup(hndls.NewMessage(message.All, middlewares.User), -10)
	dispatcher.AddHandlerToGroup(hndls.NewCallback(callbackquery.All, middlewares.User), -10)

	// handlers
	handlers := handlers.New(db, logger)
	dispatcher.AddHandlerToGroup(hndls.NewCommand("start", handlers.Start), 10)
	dispatcher.AddHandlerToGroup(hndls.NewCommand("notify", handlers.Notify), 10)
	dispatcher.AddHandlerToGroup(hndls.NewCallback(callbackquery.Prefix("guild:"), handlers.NotifyGuild), 10)
	dispatcher.AddHandlerToGroup(hndls.NewCallback(callbackquery.Prefix("channel:"), handlers.NotifyChannel), 10)

	return &bot{b, dispatcher, logger}
}

func (b *bot) Run() {
	updater := ext.NewUpdater(b.dispatcher, nil)
	if err := updater.StartPolling(b.Bot, nil); err != nil {
		b.logger.Error("failed to start polling", zap.Error(err))
		return
	}

	b.logger.Info("Bot is now running")
}
