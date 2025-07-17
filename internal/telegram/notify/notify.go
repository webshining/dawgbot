package notify

import "bot/internal/telegram/app"

type notify struct {
	app *app.AppContext
}

func New(app *app.AppContext) *notify {
	return &notify{app}
}
