package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/webshining/internal/common/database"
	"github.com/webshining/internal/common/rabbit"
	"github.com/webshining/internal/discord/handlers"
)

func main() {
	discord, err := discordgo.New("Bot " + "")
	if err != nil {
		fmt.Println("error creating Discord session:", err)
		return
	}
	defer discord.Close()

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

	handlers := handlers.New(db, ch)

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
		fmt.Println("error opening connection:", err)
		return
	}

	fmt.Println("Бот запущен. Для выхода нажмите Ctrl+C.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	<-sc
}
