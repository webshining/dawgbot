package main

import (
	"bot/internal/discord"
)

func main() {
	bot, err := discord.New()
	if err != nil {
		panic(err)
	}
	bot.Run()
}
