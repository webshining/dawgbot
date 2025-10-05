package app

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TelegramBot interface {
	SendPhoto(chatId int64, photo gotgbot.InputFileOrString, opts *gotgbot.SendPhotoOpts)
	SendMessage(chatId int64, text string, opts *gotgbot.SendMessageOpts)
}

type AppContext struct {
	Session     *discordgo.Session
	DB          *gorm.DB
	TelegramBot TelegramBot
	Logger      *zap.Logger
}

func New(session *discordgo.Session, db *gorm.DB, telegramBot TelegramBot, logger *zap.Logger) *AppContext {
	return &AppContext{
		Session:     session,
		DB:          db,
		TelegramBot: telegramBot,
		Logger:      logger,
	}
}
