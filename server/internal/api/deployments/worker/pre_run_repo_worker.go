package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/toufiq-austcse/deployit/internal/api/deployments/model"

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

type PreRunRepoWorker struct {
	config            amqp.Config
	deploymentService *service.DeploymentService
	dockerService     *service.DockerService
	eventService      *service.EventService
}

func NewPreRunRepoWorker(
	deploymentService *service.DeploymentService,
	eventService *service.EventService,
) *PreRunRepoWorker {
	return &PreRunRepoWorker{
		config: rabbit_mq.New(deployItConfig.AppConfig.RABBIT_MQ_CONFIG.EXCHANGE,
			"topic",
			deployItConfig.AppConfig.RABBIT_MQ_CONFIG.RABBIT_MQ_REPOSITORY_PRE_RUN_QUEUE,
			deployItConfig.AppConfig.RABBIT_MQ_CONFIG.RABBIT_MQ_REPOSITORY_PRE_RUN_ROUTING_KEY),
		deploymentService: deploymentService,
		eventService:      eventService,
	}
}

func (worker *PreRunRepoWorker) InitPreRunRepoSubscriber() {
	subscriber, err := amqp.NewSubscriber(worker.config, watermill.NewStdLogger(false, false))
	if err != nil {
		fmt.Println("error in pre run repo subscriber ", err.Error())
		return
	}
	fmt.Println("PreRunRepoWorker Initialized")
	messages, err := subscriber.Subscribe(
		context.Background(),
		deployItConfig.AppConfig.RABBIT_MQ_CONFIG.RABBIT_MQ_REPOSITORY_PRE_RUN_QUEUE,
	)
	if err != nil {
		fmt.Println("error in pre run repo subscriber ", err.Error())
		return
	}
	go worker.ProcessPreRunRepoMessages(messages)
}

func (worker *PreRunRepoWorker) ProcessPreRunRepoMessages(messages <-chan *message.Message) {
	for msg := range messages {
		deploymentId, event, err := worker.ProcessPreRunRepoMessage(msg)
		if err != nil {
			fmt.Println("error in processing pre run repo message ", err.Error())
			worker.eventService.WriteToFile("error in identifying port: "+err.Error(), event)

			if deploymentId != "" {
				_, updateErr := worker.deploymentService.UpdateLatestStatus(
					deploymentId,
					enums.FAILED,
					event,
					context.Background(),
				)
				if updateErr != nil {
					fmt.Println("error in updating deployment status ", updateErr.Error())
				}
			}

		}
	}
}

func (worker *PreRunRepoWorker) ProcessPreRunRepoMessage(msg *message.Message) (string, *model.Event, error) {
	defer msg.Ack()

	consumedPayload := payloads.PreRunRepoWorkerPayload{}
	if err := json.Unmarshal(msg.Payload, &consumedPayload); err != nil {
		return "", nil, err
	}
	fmt.Println("consumed pre run job ", consumedPayload)

	deployment := worker.deploymentService.FindById(
		consumedPayload.DeploymentId,
		context.Background(),
	)
	if deployment == nil {
		return consumedPayload.DeploymentId, nil, app_errors.DeploymentNotFoundError
	}
	event, err := worker.eventService.FindById(consumedPayload.EventId)
	if err != nil {
		fmt.Println("error in finding event ", err.Error())
	}
	if deployment.DockerImageTag == nil {
		return consumedPayload.DeploymentId, event, app_errors.DockerImageTagNotFoundError
	}
	if deployment.ContainerId != nil {
		if removeErr := worker.dockerService.RemoveContainer(*deployment.ContainerId); removeErr != nil {
			return consumedPayload.DeploymentId, event, removeErr
		}

		fmt.Println("container removed ", deployment.ContainerId)
		if _, updateErr := worker.deploymentService.UpdateDeployment(consumedPayload.DeploymentId, map[string]interface{}{
			"container_id": nil,
		}, event, context.Background()); updateErr != nil {
			return consumedPayload.DeploymentId, event, updateErr
		}

	}
	containerId, preRunErr := worker.dockerService.RunContainer(*deployment.DockerImageTag,
		deployment.Env, nil)
	if preRunErr != nil {
		return consumedPayload.DeploymentId, event, preRunErr
	}

	fmt.Println("docker image pre run successfully...")
	worker.eventService.WriteToFile("identifying port", event)

	_, updateErr := worker.deploymentService.UpdateDeployment(
		consumedPayload.DeploymentId,
		map[string]interface{}{
			"latest_status": enums.DEPLOYING,
			"container_id":  containerId,
		},
		event,
		context.Background(),
	)

	if updateErr != nil {
		return consumedPayload.DeploymentId, event, updateErr
	}
	return consumedPayload.DeploymentId, event, nil
}

func (worker *PreRunRepoWorker) PublishPreRunRepoJob(
	preRunRepoPayload payloads.PreRunRepoWorkerPayload,
) error {
	publisher, err := amqp.NewPublisher(worker.config, watermill.NewStdLogger(false, false))
	if err != nil {
		return err
	}

	payload, _ := json.Marshal(preRunRepoPayload)
	msg := message.NewMessage(watermill.NewUUID(), payload)
	err = publisher.Publish(
		deployItConfig.AppConfig.RABBIT_MQ_CONFIG.RABBIT_MQ_REPOSITORY_PRE_RUN_QUEUE,
		msg,
	)
	if err != nil {
		return err
	}
	return publisher.Close()
}
