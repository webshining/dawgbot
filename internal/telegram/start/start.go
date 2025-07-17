package start

import "bot/internal/telegram/app"

type start struct {
	app *app.AppContext
}

func New(app *app.AppContext) *start {
	return &start{app}
}
