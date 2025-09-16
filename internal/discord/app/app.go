package app

import (
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AppContext struct {
	Session *discordgo.Session
	DB      *gorm.DB
	Logger  *zap.Logger
}

func New(session *discordgo.Session, db *gorm.DB, logger *zap.Logger) *AppContext {
	return &AppContext{
		Session: session,
		DB:      db,
		Logger:  logger,
	}
}
