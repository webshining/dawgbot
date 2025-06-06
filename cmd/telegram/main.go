package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/webshining/internal/telegram"
)

func main() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe("0.0.0.0:2112", nil)
	}()

	bot, err := telegram.New()
	if err != nil {
		panic(err)
	}
	bot.Run()
}
