package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	config2 "github.com/toufiq-austcse/deployit/config"
	"github.com/toufiq-austcse/deployit/pkg/rabbit_mq"
)

type PullRepoWorker struct {
	config amqp.Config
}

func NewPullRepoWorker() *PullRepoWorker {
	return &PullRepoWorker{config: rabbit_mq.New(config2.AppConfig.RABBIT_MQ_CONFIG.EXCHANGE,
		"topic",
		config2.AppConfig.RABBIT_MQ_CONFIG.REPOSITORY_PULL_QUEUE,
		config2.AppConfig.RABBIT_MQ_CONFIG.REPOSITORY_PULL_ROUTING_KEY)}
}

func (pullRepoWorker PullRepoWorker) InitPullRepoSubscriber() {
	subscriber, err := amqp.NewSubscriber(pullRepoWorker.config, watermill.NewStdLogger(false, false))
	if err != nil {
		fmt.Println("error in pull repo subscriber ", err.Error())
		return
	}
	fmt.Println("PullRepoWorker Initialized")
	messages, err := subscriber.Subscribe(context.Background(), config2.AppConfig.RABBIT_MQ_CONFIG.REPOSITORY_PULL_QUEUE)
	if err != nil {
		fmt.Println("error in pull repo subscriber ", err.Error())
		return
	}
	go func() {
		type Data struct {
			Name string `json:"name"`
			Age  int64  `json:"age"`
		}
		consumedPayload := Data{}
		for msg := range messages {
			err = json.Unmarshal(msg.Payload, &consumedPayload)
			if err != nil {
				fmt.Println("error in parsing ", err)
				continue
			}
			fmt.Println(msg.UUID, consumedPayload.Name, consumedPayload.Age)

			// we need to Acknowledge that we received and processed the message,
			// otherwise, it will be resent over and over again.
			msg.Ack()
		}
	}()

}

func (pullRepoWorker PullRepoWorker) PublishPullRepoJob() error {
	publisher, err := amqp.NewPublisher(pullRepoWorker.config, watermill.NewStdLogger(false, false))
	if err != nil {
		return err
	}
	data := struct {
		Name string `json:"name"`
		Age  int64  `json:"age"`
	}{
		Name: "Sadi",
		Age:  21,
	}
	payload, _ := json.Marshal(data)
	msg := message.NewMessage(watermill.NewUUID(), payload)
	err = publisher.Publish(config2.AppConfig.RABBIT_MQ_CONFIG.REPOSITORY_PULL_QUEUE, msg)
	if err != nil {
		return err
	}
	return publisher.Close()
}
