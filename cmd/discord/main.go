package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/webshining/internal/common/database"
	"github.com/webshining/internal/common/rabbit"
	"github.com/webshining/internal/discord/handlers"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	godotenv.Load()

	discord, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		logger.Error("error creating Discord session", zap.Error(err))
		return
	}
	defer discord.Close()

	db, err := database.New(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC", os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT")))
	if err != nil {
		logger.Error("error connecting to database", zap.Error(err))
		return
	}

	conn, ch, err := rabbit.New(fmt.Sprintf("amqp://%s:%s@%s:%s/", os.Getenv("RB_USER"), os.Getenv("RB_PASS"), os.Getenv("RB_HOST"), os.Getenv("RB_PORT")))
	if err != nil {
		logger.Error("error connecting to RabbitMQ", zap.Error(err))
		return
	}
	defer conn.Close()
	defer ch.Close()

	handlers := handlers.New(db, ch, logger)

	discord.Identify.Intents = discordgo.IntentsGuildVoiceStates | discordgo.IntentsGuilds

	// voice handlers
	discord.AddHandler(handlers.VoiceJoinHandler)
	// guild handlers
	discord.AddHandler(handlers.GuildAddHandler)
	discord.AddHandler(handlers.GuildUpdateHandler)
	discord.AddHandler(handlers.GuildDeleteHandler)
	// channel handlers
	discord.AddHandler(handlers.ChannelAddHandler)
	discord.AddHandler(handlers.ChannelUpdateHandler)
	discord.AddHandler(handlers.ChannelDeleteHandler)

	if err := discord.Open(); err != nil {
		logger.Error("error opening connection to Discord", zap.Error(err))
		return
	}

	logger.Info("Bot is now running. Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc
}
