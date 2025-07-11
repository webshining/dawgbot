package middlewares

import (
	"bot/internal/telegram/app"

	"gorm.io/gorm"
)

type middlewares struct {
	db *gorm.DB
}

func New(app *app.AppContext) *middlewares {
	return &middlewares{db: app.DB}
}
