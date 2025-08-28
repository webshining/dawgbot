package app

import (
	"bot/internal/common/broker"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AppContext struct {
	Bot    *gotgbot.Bot
	DB     *gorm.DB
	Logger *zap.Logger
	Broker *broker.Broker
}

func New(bot *gotgbot.Bot, db *gorm.DB, broker *broker.Broker, logger *zap.Logger) *AppContext {
	return &AppContext{
		Bot:    bot,
		DB:     db,
		Logger: logger,
		Broker: broker,
	}
}
