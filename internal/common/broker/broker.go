package broker

import (
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

type Broker struct {
	Client mqtt.Client
}

func New(id string, logger *zap.Logger) *Broker {
	opts := mqtt.NewClientOptions().AddBroker(os.Getenv("BROKER_ADDR")).SetClientID(id)
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		logger.Error("error connection to mqtt", zap.Error(token.Error()))
	}
	return &Broker{c}
}

func (b *Broker) Publish(topic string, payload interface{}) mqtt.Token {
	return b.Client.Publish(topic, 0, false, payload)
}

func (b *Broker) Subscribe(topic string, callback mqtt.MessageHandler) mqtt.Token {
	return b.Client.Subscribe(topic, 0, callback)
}
