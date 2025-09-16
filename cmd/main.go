package main

import (
	"bot/internal/discord"
	"bot/internal/telegram"
)

func main() {
	telegramBot, err := telegram.New()
	if err != nil {
		panic(err)
	}

	discordBot, err := discord.New()
	if err != nil {
		panic(err)
	}

	go telegramBot.Run()
	go discordBot.Run()

	select {}
}
