package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/toufiq-austcse/deployit/pkg/utils"

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

type RunRepoWorker struct {
	config            amqp.Config
	deploymentService *service.DeploymentService
	dockerService     *service.DockerService
	eventService      *service.EventService
}

func NewRunRepoWorker(deploymentService *service.DeploymentService, eventService *service.EventService) *RunRepoWorker {
	return &RunRepoWorker{
		config: rabbit_mq.New(deployItConfig.AppConfig.RABBIT_MQ_CONFIG.EXCHANGE,
			"topic",
			deployItConfig.AppConfig.RABBIT_MQ_CONFIG.REPOSITORY_RUN_QUEUE,
			deployItConfig.AppConfig.RABBIT_MQ_CONFIG.REPOSITORY_RUN_ROUTING_KEY),
		deploymentService: deploymentService,
		eventService:      eventService,
	}
}

func (worker *RunRepoWorker) InitRunRepoSubscriber() {
	subscriber, err := amqp.NewSubscriber(worker.config, watermill.NewStdLogger(false, false))
	if err != nil {
		fmt.Println("error in run repo subscriber ", err.Error())
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
	go worker.ProcessRunRepoMessages(messages)
}

func (worker *RunRepoWorker) ProcessRunRepoMessages(messages <-chan *message.Message) {
	for msg := range messages {
		deploymentId, lastDeploymentInitiateAt, event, err := worker.ProcessRunRepoMessage(msg)
		if err != nil {
			if deploymentId != "" {
				fmt.Println("error in processing run repo message ", err.Error())
				utils.WriteToFile("error in running docker container: "+err.Error(), event)
				if err.Error() == app_errors.ContainerPortNotFoundError.Error() &&
					lastDeploymentInitiateAt != nil {
					timeElapsed := time.Since(*lastDeploymentInitiateAt)
					if int(
						timeElapsed.Minutes(),
					) < deployItConfig.AppConfig.MAX_DEPLOYING_STATUS_TIME_IN_MINUTES {
						fmt.Println("time elapsed ", timeElapsed.Minutes())
						continue
					}
				}
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

func (worker *RunRepoWorker) ProcessRunRepoMessage(
	msg *message.Message,
) (string, *time.Time, *model.Event, error) {
	defer msg.Ack()

	consumedPayload := payloads.RunRepoWorkerPayload{}
	if err := json.Unmarshal(msg.Payload, &consumedPayload); err != nil {
		return "", nil, nil, err
	}
	fmt.Println("consumed run job ", consumedPayload)

	deployment := worker.deploymentService.FindById(
		consumedPayload.DeploymentId,
		context.Background(),
	)
	if deployment == nil {
		return consumedPayload.DeploymentId, nil, nil, app_errors.DeploymentNotFoundError
	}

	existingEvent, err := worker.eventService.FindById(consumedPayload.EventId)
	if err != nil {
		fmt.Println("error in finding event ", err.Error())
	}

	if deployment.DockerImageTag == nil {
		return consumedPayload.DeploymentId, deployment.LastDeploymentInitiatedAt, existingEvent, app_errors.DockerImageTagNotFoundError
	}
	if deployment.ContainerId == nil {
		return consumedPayload.DeploymentId, deployment.LastDeploymentInitiatedAt, existingEvent, app_errors.ContainerPortNotFoundError
	}
	port, portErr := worker.dockerService.GetTcpPort(*deployment.ContainerId)
	if portErr != nil {
		return consumedPayload.DeploymentId, deployment.LastDeploymentInitiatedAt, existingEvent, app_errors.ContainerPortNotFoundError
	}
	fmt.Println("port found ", *port)
	utils.WriteToFile("Detected service running port "+*port, existingEvent)
	if removeErr := worker.dockerService.RemoveContainer(*deployment.ContainerId); removeErr != nil {
		return consumedPayload.DeploymentId, deployment.LastDeploymentInitiatedAt, existingEvent, removeErr
	}

	fmt.Println("container removed ", deployment.ContainerId)
	if _, updateErr := worker.deploymentService.UpdateDeployment(consumedPayload.DeploymentId, map[string]interface{}{
		"container_id": nil,
	}, existingEvent, context.Background()); updateErr != nil {
		return consumedPayload.DeploymentId, deployment.LastDeploymentInitiatedAt, existingEvent, updateErr
	}

	utils.WriteToFile("running your service", existingEvent)
	containerId, runErr := worker.dockerService.RunContainer(
		*deployment.DockerImageTag,
		deployment.Env,
		port,
	)
	if runErr != nil {
		return consumedPayload.DeploymentId, deployment.LastDeploymentInitiatedAt, existingEvent, runErr
	}
	fmt.Println("docker image run successfully...")
	utils.WriteToFile("deployed successfully", existingEvent)

	_, updateErr := worker.deploymentService.UpdateDeployment(
		consumedPayload.DeploymentId,
		map[string]interface{}{
			"latest_status":    enums.LIVE,
			"last_deployed_at": time.Now(),
			"container_id":     containerId,
		},
		existingEvent,
		context.Background(),
	)

	if updateErr != nil {
		return consumedPayload.DeploymentId, nil, existingEvent, updateErr
	}
	return consumedPayload.DeploymentId, deployment.LastDeploymentInitiatedAt, existingEvent, nil
}

func (worker *RunRepoWorker) PublishRunRepoJob(runRepoPayload payloads.RunRepoWorkerPayload) error {
	publisher, err := amqp.NewPublisher(worker.config, watermill.NewStdLogger(false, false))
	if err != nil {
		return err
	}

	payload, _ := json.Marshal(runRepoPayload)
	msg := message.NewMessage(watermill.NewUUID(), payload)
	err = publisher.Publish(deployItConfig.AppConfig.RABBIT_MQ_CONFIG.REPOSITORY_RUN_QUEUE, msg)
	if err != nil {
		return err
	}
	return publisher.Close()
}
