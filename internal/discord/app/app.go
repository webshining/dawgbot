package app

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AppContext struct {
	Session     *discordgo.Session
	DB          *gorm.DB
	TelegramBot *gotgbot.Bot
	Logger      *zap.Logger
}

func New(session *discordgo.Session, db *gorm.DB, telegramBot *gotgbot.Bot, logger *zap.Logger) *AppContext {
	return &AppContext{
		Session:     session,
		DB:          db,
		TelegramBot: telegramBot,
		Logger:      logger,
	}
}
