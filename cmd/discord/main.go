package main

import (
	"bot/internal/discord"
)

func main() {
	bot := discord.New()
	bot.Run()
}
