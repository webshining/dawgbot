package main

import "github.com/webshining/internal/discord"

func main() {
	bot, err := discord.New()
	if err != nil {
		panic(err)
	}
	bot.Run()
}
