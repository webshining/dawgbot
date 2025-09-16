package app

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AppContext struct {
	Bot    *gotgbot.Bot
	DB     *gorm.DB
	Logger *zap.Logger
}

func New(bot *gotgbot.Bot, db *gorm.DB, logger *zap.Logger) *AppContext {
	return &AppContext{
		Bot:    bot,
		DB:     db,
		Logger: logger,
	}
}
