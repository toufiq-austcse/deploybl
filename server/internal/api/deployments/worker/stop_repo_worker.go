package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	deployItConfig "github.com/toufiq-austcse/deployit/config"
	"github.com/toufiq-austcse/deployit/enums"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/service"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/worker/payloads"
	"github.com/toufiq-austcse/deployit/pkg/app_errors"
	"github.com/toufiq-austcse/deployit/pkg/rabbit_mq"
)

type StopRepoWorker struct {
	config            amqp.Config
	deploymentService *service.DeploymentService
	dockerService     *service.DockerService
}

func NewStopRepoWorker(deploymentService *service.DeploymentService) *StopRepoWorker {
	return &StopRepoWorker{
		config: rabbit_mq.New(deployItConfig.AppConfig.RABBIT_MQ_CONFIG.EXCHANGE,
			"topic",
			deployItConfig.AppConfig.RABBIT_MQ_CONFIG.REPOSITORY_STOP_QUEUE,
			deployItConfig.AppConfig.RABBIT_MQ_CONFIG.REPOSITORY_STOP_ROUTING_KEY),
		deploymentService: deploymentService,
	}
}

func (worker *StopRepoWorker) InitStopRepoSubscriber() {
	subscriber, err := amqp.NewSubscriber(worker.config, watermill.NewStdLogger(false, false))
	if err != nil {
		fmt.Println("error in stop repo subscriber ", err.Error())
		return
	}
	fmt.Println("RunRepoSubscriber Initialized")
	messages, err := subscriber.Subscribe(
		context.Background(),
		deployItConfig.AppConfig.RABBIT_MQ_CONFIG.REPOSITORY_PULL_QUEUE,
	)
	if err != nil {
		fmt.Println("error in run repo subscriber ", err.Error())
		return
	}
	go worker.ProcessStopRepoMessages(messages)
}

func (worker *StopRepoWorker) ProcessStopRepoMessages(messages <-chan *message.Message) {
	for msg := range messages {
		deploymentId, err := worker.ProcessStopRepoMessage(msg)
		if err != nil {
			fmt.Println("error in processing run repo message ", err.Error())

			if deploymentId != "" {
				_, _, updateErr := worker.deploymentService.UpdateLatestStatus(
					deploymentId,
					enums.FAILED,
					"",
					context.Background(),
				)
				if updateErr != nil {
					fmt.Println("error in updating deployment status ", updateErr.Error())
				}
			}

		}
	}
}

func (worker *StopRepoWorker) ProcessStopRepoMessage(msg *message.Message) (string, error) {
	defer msg.Ack()

	consumedPayload := payloads.StopRepoWorkerPayload{}
	if err := json.Unmarshal(msg.Payload, &consumedPayload); err != nil {
		return "", err
	}
	fmt.Println("consumed run job ", consumedPayload)

	deployment := worker.deploymentService.FindById(
		consumedPayload.DeploymentId,
		context.Background(),
	)
	if deployment == nil {
		return consumedPayload.DeploymentId, app_errors.DeploymentNotFoundError
	}
	if deployment.ContainerId == nil {
		return consumedPayload.DeploymentId, app_errors.ContainerNotFoundError
	}
	if !worker.deploymentService.IsStopAble(deployment) {
		return consumedPayload.DeploymentId, app_errors.DeploymentNotStoppableError
	}
	if err := worker.dockerService.StopContainer(*deployment.ContainerId); err != nil {
		return consumedPayload.DeploymentId, err
	}

	_, _, err := worker.deploymentService.UpdateLatestStatus(
		consumedPayload.DeploymentId,
		enums.STOPPED,
		"",
		context.Background(),
	)
	if err != nil {
		return consumedPayload.DeploymentId, err
	}
	return consumedPayload.DeploymentId, nil
}

func (worker *StopRepoWorker) PublishStopRepoJob(
	stopRepoPayload payloads.StopRepoWorkerPayload,
) error {
	publisher, err := amqp.NewPublisher(worker.config, watermill.NewStdLogger(false, false))
	if err != nil {
		return err
	}

	payload, _ := json.Marshal(stopRepoPayload)
	msg := message.NewMessage(watermill.NewUUID(), payload)
	err = publisher.Publish(deployItConfig.AppConfig.RABBIT_MQ_CONFIG.REPOSITORY_RUN_QUEUE, msg)
	if err != nil {
		return err
	}
	return publisher.Close()
}
