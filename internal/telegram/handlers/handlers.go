package handlers

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/webshining/internal/telegram/app"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type handlers struct {
	db     *gorm.DB
	logger *zap.Logger
	rabbit *amqp.Channel
}

func New(app *app.AppContext) *handlers {
	rabbit, _ := app.Rabbit.Channel()
	rabbit.Consume("playlist", "", true, false, false, false, nil)
	return &handlers{
		db:     app.DB,
		logger: app.Logger,
		rabbit: rabbit,
	}
}
