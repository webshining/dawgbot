package telegram

import (
	"context"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	hndls "github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"bot/internal/config"
	"bot/internal/telegram/handlers"
	"bot/internal/telegram/middlewares"
)

type Bot struct {
	bot        *gotgbot.Bot
	dispatcher *ext.Dispatcher
	logger     *zap.Logger
}

func New(config config.TelegramConfig, database *gorm.DB, logger *zap.Logger) *Bot {
	// setup new bot session
	b, err := gotgbot.NewBot(config.Token, nil)
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
	middlewares := middlewares.New(database, logger)
	dispatcher.AddHandlerToGroup(hndls.NewMessage(message.All, middlewares.User), -10)
	dispatcher.AddHandlerToGroup(hndls.NewCallback(callbackquery.All, middlewares.User), -10)

	// handlers
	handlers := handlers.New(database, logger)
	dispatcher.AddHandlerToGroup(hndls.NewCommand("start", handlers.Start), 10)
	dispatcher.AddHandlerToGroup(hndls.NewCommand("notify", handlers.Notify), 10)
	dispatcher.AddHandlerToGroup(hndls.NewCallback(callbackquery.Prefix("guild:"), handlers.NotifyGuild), 10)
	dispatcher.AddHandlerToGroup(hndls.NewCallback(callbackquery.Prefix("channel:"), handlers.NotifyChannel), 10)

	return &Bot{b, dispatcher, logger}
}

func (b *Bot) Run(ctx context.Context) {
	defer b.bot.Close(nil)

	updater := ext.NewUpdater(b.dispatcher, nil)
	if err := updater.StartPolling(b.bot, nil); err != nil {
		b.logger.Error("failed to start polling", zap.Error(err))
		return
	}

	b.logger.Info("Bot is now running")
	defer b.logger.Info("Bot is shuted down")

	<-ctx.Done()
}
