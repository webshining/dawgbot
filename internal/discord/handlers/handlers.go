package handlers

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Handlers struct {
	DB     *gorm.DB
	AMQP   *amqp.Channel
	Logger *zap.Logger
}

func New(db *gorm.DB, amqp *amqp.Channel, logger *zap.Logger) *Handlers {
	return &Handlers{DB: db, AMQP: amqp, Logger: logger}
}
