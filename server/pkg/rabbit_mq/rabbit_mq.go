package rabbit_mq

import (
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	config2 "github.com/toufiq-austcse/deployit/config"
)

func New(exchangeName, exchangeType, queueName, routingKey string) amqp.Config {
	queueGen := amqp.GenerateQueueNameConstant(queueName)
	pubSubConfig := amqp.NewDurablePubSubConfig(config2.AppConfig.RABBIT_MQ_CONFIG.URL, queueGen)
	pubSubConfig.Exchange = amqp.ExchangeConfig{
		GenerateName: func(topic string) string {
			return exchangeName
		},
		Type: exchangeType,
	}
	pubSubConfig.QueueBind = amqp.QueueBindConfig{
		GenerateRoutingKey: func(topic string) string {
			return routingKey
		},
	}
	pubSubConfig.Publish = amqp.PublishConfig{
		GenerateRoutingKey: func(topic string) string {
			return routingKey
		},
		ChannelPoolSize: 5,
	}
	return pubSubConfig
}
