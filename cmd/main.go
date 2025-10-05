package main

import (
	"bot/internal/config"
	"bot/internal/database"
	"bot/internal/discord"
	"bot/internal/logger"
	"bot/internal/telegram"
	"context"
	"os"
	"os/signal"
	"sync"
)

func main() {
	logger := logger.New()
	defer logger.Sync()
	config := config.MustLoad(logger)

	database := database.MustConnect(config.Database, logger)

	telegramBot := telegram.New(config.Telegram, database, logger)
	discordBot := discord.New(&config.Discord, database, telegramBot, logger)

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	defer func() {
		cancel()
		wg.Wait()
	}()

	wg.Add(2)
	go func() {
		defer wg.Done()
		telegramBot.Run(ctx)
	}()
	go func() {
		defer wg.Done()
		discordBot.Run(ctx)
	}()

	logger.Info("Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc
}
