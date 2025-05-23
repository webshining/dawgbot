package main

import "github.com/webshining/internal/telegram"

func main() {
	bot := telegram.New()
	if bot != nil {
		bot.Run()
	}
}
