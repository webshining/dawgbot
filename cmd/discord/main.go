package main

import "github.com/webshining/internal/discord"

func main() {
	bot := discord.New()
	if bot != nil {
		bot.Run()
	}
}
