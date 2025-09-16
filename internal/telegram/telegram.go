package telegram

import (
	"fmt"
	"os"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	hndls "github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"bot/internal/common/database"
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

	// setup new database connection
	db, err := database.New(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT")))
	if err != nil {
		logger.Fatal("error connecting to database", zap.Error(err))
	}

	// setup new bot session
	b, err := gotgbot.NewBot(os.Getenv("TELEGRAM_BOT_TOKEN"), nil)
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
