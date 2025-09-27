package broker

import (
	"bot/internal/config"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

type Broker struct {
	Client mqtt.Client
}

func Connect(id string, config config.BrokerConfig, logger *zap.Logger) (*Broker, error) {
	opts := mqtt.NewClientOptions().AddBroker(config.Addr).SetClientID(id)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return &Broker{client}, nil
}

func MustConnect(id string, config config.BrokerConfig, logger *zap.Logger) *Broker {
	client, err := Connect(id, config, logger)
	if err != nil {
		logger.Named("Broker").Fatal("failed to connect to emqx", zap.Error(err))
	}
	return client
}

func (b *Broker) Publish(topic string, payload interface{}) mqtt.Token {
	println("message sended")
	return b.Client.Publish(topic, 0, false, payload)
}

func (b *Broker) Subscribe(topic string, callback mqtt.MessageHandler) mqtt.Token {
	return b.Client.Subscribe(topic, 0, callback)
}
