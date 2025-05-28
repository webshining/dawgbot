package main

import "github.com/webshining/internal/telegram"

func main() {
	bot, err := telegram.New()
	if err != nil {
		panic(err)
	}
	bot.Run()
}
