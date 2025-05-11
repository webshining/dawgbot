package main

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
	"go.uber.org/zap"

	"github.com/webshining/internal/common/database"
	"github.com/webshining/internal/common/rabbit"
	hndls "github.com/webshining/internal/telegram/handlers"
	"github.com/webshining/internal/telegram/middlewares"
	"github.com/webshining/internal/telegram/notifier"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	godotenv.Load()

	b, err := gotgbot.NewBot(os.Getenv("TELEGRAM_BOT_TOKEN"), nil)
	if err != nil {
		fmt.Println("failed to create new bot:", err)
		return
	}

	db, err := database.New(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT")))
	if err != nil {
		fmt.Println("failed to connect database:", err)
		return
	}

	conn, ch, err := rabbit.New(fmt.Sprintf("amqp://%s:%s@%s:%s/", os.Getenv("RB_USER"), os.Getenv("RB_PASS"), os.Getenv("RB_HOST"), os.Getenv("RB_PORT")))
	if err != nil {
		fmt.Println("rabbit error:", err)
	}
	defer conn.Close()
	defer ch.Close()

	dispatcher := ext.NewDispatcher(nil)
	updater := ext.NewUpdater(dispatcher, nil)

	middlewares := middlewares.New(db)
	hndl := hndls.New(db, logger)

	dispatcher.AddHandlerToGroup(handlers.NewMessage(message.Command, middlewares.UserMiddleware), -10)
	dispatcher.AddHandlerToGroup(handlers.NewCallback(callbackquery.All, middlewares.UserMiddleware), -10)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("start", hndl.StartHandler), 10)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("notify", hndl.NotifyHandler), 10)
	dispatcher.AddHandlerToGroup(handlers.NewCallback(callbackquery.Prefix("guild:"), hndl.NotifyGuildHandler), 10)
	dispatcher.AddHandlerToGroup(handlers.NewCallback(callbackquery.Prefix("channel:"), hndl.NotifyChannelHandler), 10)

	err = updater.StartPolling(b, nil)
	if err != nil {
		logger.Error("failed to start polling", zap.Error(err))
	}

	notifier := notifier.New(ch, b, db, logger)

	go func() {
		err := notifier.Start()
		if err != nil {
			logger.Error("failed to start notifier", zap.Error(err))
		}
	}()

	logger.Info("Bot is now running. Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc
}
