package app

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AppContext struct {
	Bot    *gotgbot.Bot
	Rabbit *amqp.Connection
	DB     *gorm.DB
	Logger *zap.Logger
}

func New(bot *gotgbot.Bot, rabbit *amqp.Connection, db *gorm.DB, logger *zap.Logger) *AppContext {
	return &AppContext{
		Bot:    bot,
		Rabbit: rabbit,
		DB:     db,
		Logger: logger,
	}
}
