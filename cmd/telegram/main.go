package main

import (
	"bot/internal/telegram"
)

func main() {
	bot := telegram.New()
	bot.Run()
}
