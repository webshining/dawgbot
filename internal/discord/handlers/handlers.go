package handlers

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type Handlers struct {
	DB   *gorm.DB
	AMQP *amqp.Channel
}

func New(db *gorm.DB, amqp *amqp.Channel) *Handlers {
	return &Handlers{DB: db, AMQP: amqp}
}
