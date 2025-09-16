package main

import (
	"bot/internal/discord"
	"bot/internal/telegram"
)

func main() {
	telegramBot := telegram.New()

	discordBot := discord.New(telegramBot.Bot)

	go telegramBot.Run()
	go discordBot.Run()

	select {}
}
