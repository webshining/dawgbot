package main

import (
	"bot/internal/telegram"
)

func main() {
	bot, err := telegram.New()
	if err != nil {
		panic(err)
	}
	bot.Run()
}
