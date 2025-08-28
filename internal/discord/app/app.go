package app

import (
	"bot/internal/common/broker"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AppContext struct {
	Session *discordgo.Session
	DB      *gorm.DB
	Logger  *zap.Logger
	Broker  *broker.Broker
}

func New(session *discordgo.Session, db *gorm.DB, broker *broker.Broker, logger *zap.Logger) *AppContext {
	return &AppContext{
		Session: session,
		DB:      db,
		Logger:  logger,
		Broker:  broker,
	}
}
