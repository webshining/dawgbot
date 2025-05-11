package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"

	"github.com/webshining/internal/common/database"
	"github.com/webshining/internal/common/rabbit"
	hndls "github.com/webshining/internal/telegram/handlers"
	"github.com/webshining/internal/telegram/middlewares"
	"github.com/webshining/internal/telegram/notifier"
)

func main() {
	db, err := database.New()
	if err != nil {
		fmt.Println("failed to connect database:", err)
		return
	}

	conn, ch, err := rabbit.New()
	if err != nil {
		fmt.Println("rabbit error:", err)
	}
	defer conn.Close()
	defer ch.Close()

	b, err := gotgbot.NewBot("", nil)
	if err != nil {
		fmt.Println("failed to create new bot:", err)
		return
	}

	dispatcher := ext.NewDispatcher(nil)
	updater := ext.NewUpdater(dispatcher, nil)

	middlewares := middlewares.New(db)
	hndl := hndls.New(db)

	dispatcher.AddHandlerToGroup(handlers.NewMessage(message.Command, middlewares.UserMiddleware), -10)
	dispatcher.AddHandlerToGroup(handlers.NewCallback(callbackquery.All, middlewares.UserMiddleware), -10)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("start", hndl.StartHandler), 10)
	dispatcher.AddHandlerToGroup(handlers.NewCommand("notify", hndl.NotifyHandler), 10)
	dispatcher.AddHandlerToGroup(handlers.NewCallback(callbackquery.Prefix("guild:"), hndl.NotifyGuildHandler), 10)
	dispatcher.AddHandlerToGroup(handlers.NewCallback(callbackquery.Prefix("channel:"), hndl.NotifyChannelHandler), 10)

	err = updater.StartPolling(b, nil)
	if err != nil {
		fmt.Println("failed to start polling:", err)
	}

	notifier := notifier.New(ch, b, db)

	go func() {
		err := notifier.Start()
		if err != nil {
			log.Printf("Ошибка при запуске консьюмера: %v", err)
		}
	}()

	fmt.Println("Бот запущен. Для выхода нажмите Ctrl+C.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc
}
