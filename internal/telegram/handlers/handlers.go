package handlers

import (
	"github.com/webshining/internal/telegram/app"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type handlers struct {
	db     *gorm.DB
	logger *zap.Logger
}

func New(app *app.AppContext) *handlers {
	return &handlers{
		db:     app.DB,
		logger: app.Logger,
	}
}
