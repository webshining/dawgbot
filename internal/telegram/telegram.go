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
	"github.com/webshining/internal/common/database"
	"github.com/webshining/internal/common/rabbit"
	hndls "github.com/webshining/internal/telegram/handlers"
	"github.com/webshining/internal/telegram/middlewares"
	"github.com/webshining/internal/telegram/notifier"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Bot struct {
	Bot        *gotgbot.Bot
	Dispatcher *ext.Dispatcher
	DB         *gorm.DB
	Rabbit     *amqp.Channel
	logger     *zap.Logger
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
	b, err := gotgbot.NewBot(os.Getenv("TELEGRAM_BOT_TOKEN"), nil)
	if err != nil {
		logger.Error("failed to create new bot:", zap.Error(err))
		return nil
	}
	dispatcher := ext.NewDispatcher(nil)

	// setup bot handlers
	middlewares := middlewares.New(db)
	hndl := hndls.New(db, logger)

	dispatcher.AddHandlerToGroup(handlers.NewMessage(message.Command, middlewares.UserMiddleware), -10)
	dispatcher.AddHandlerToGroup(handlers.NewCallback(callbackquery.All, middlewares.UserMiddleware), -10)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("start", hndl.StartHandler), 10)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("notify", hndl.NotifyHandler), 10)
	dispatcher.AddHandlerToGroup(handlers.NewCallback(callbackquery.Prefix("guild:"), hndl.NotifyGuildHandler), 10)
	dispatcher.AddHandlerToGroup(handlers.NewCallback(callbackquery.Prefix("channel:"), hndl.NotifyChannelHandler), 10)

	return &Bot{
		Bot:        b,
		Dispatcher: dispatcher,
		DB:         db,
		Rabbit:     amqp_ch,
		logger:     logger,
	}
}

func (b *Bot) Run() {
	updater := ext.NewUpdater(b.Dispatcher, nil)
	if err := updater.StartPolling(b.Bot, nil); err != nil {
		b.logger.Error("failed to start polling", zap.Error(err))
	}

	notifier := notifier.New(b.Rabbit, b.Bot, b.DB, b.logger)

	go func() {
		if err := notifier.Start(); err != nil {
			b.logger.Error("failed to start notifier", zap.Error(err))
		}
	}()

	b.logger.Info("Bot is now running. Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc
}
