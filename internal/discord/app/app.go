package app

import (
	"github.com/bwmarrin/discordgo"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AppContext struct {
	Session *discordgo.Session
	DB      *gorm.DB
	Rabbit  *amqp.Connection
	Logger  *zap.Logger
}

func New(session *discordgo.Session, db *gorm.DB, rabbit *amqp.Connection, logger *zap.Logger) *AppContext {
	return &AppContext{
		Session: session,
		DB:      db,
		Rabbit:  rabbit,
		Logger:  logger,
	}
}
