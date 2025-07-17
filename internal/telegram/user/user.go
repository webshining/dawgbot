package user

import "bot/internal/telegram/app"

type user struct {
	app *app.AppContext
}

func New(app *app.AppContext) *user {
	return &user{app}
}
