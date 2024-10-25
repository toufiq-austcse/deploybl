package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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
}

func NewRunRepoWorker(deploymentService *service.DeploymentService) *RunRepoWorker {
	return &RunRepoWorker{
		config: rabbit_mq.New(deployItConfig.AppConfig.RABBIT_MQ_CONFIG.EXCHANGE,
			"topic",
			deployItConfig.AppConfig.RABBIT_MQ_CONFIG.REPOSITORY_RUN_QUEUE,
			deployItConfig.AppConfig.RABBIT_MQ_CONFIG.REPOSITORY_RUN_ROUTING_KEY),
		deploymentService: deploymentService,
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
		deploymentId, lastDeploymentInitiateAt, err := worker.ProcessRunRepoMessage(msg)
		if err != nil {
			if deploymentId != "" {
				if err.Error() == app_errors.ContainerPortNotFoundError.Error() &&
					lastDeploymentInitiateAt != nil {
					timeElapsed := time.Since(*lastDeploymentInitiateAt)
					if int(
						timeElapsed.Minutes(),
					) < deployItConfig.AppConfig.MAX_DEPLOYING_STATUS_TIME_IN_MINUTES {
						continue
					}
				}
				_, updateErr := worker.deploymentService.UpdateLatestStatus(
					deploymentId,
					enums.FAILED,
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
) (string, *time.Time, error) {
	defer msg.Ack()

	consumedPayload := payloads.RunRepoWorkerPayload{}
	if err := json.Unmarshal(msg.Payload, &consumedPayload); err != nil {
		return "", nil, err
	}
	fmt.Println("consumed run job ", consumedPayload)

	deployment := worker.deploymentService.FindById(
		consumedPayload.DeploymentId,
		context.Background(),
	)
	if deployment == nil {
		return consumedPayload.DeploymentId, nil, app_errors.DeploymentNotFoundError
	}
	if deployment.DockerImageTag == nil {
		return consumedPayload.DeploymentId, deployment.LastDeploymentInitiatedAt, app_errors.DockerImageTagNotFoundError
	}
	port, portErr := worker.dockerService.GetTcpPort(*deployment.ContainerId)
	if portErr != nil {
		return consumedPayload.DeploymentId, deployment.LastDeploymentInitiatedAt, app_errors.ContainerPortNotFoundError
	}
	fmt.Println("port found ", *port)
	if removeErr := worker.dockerService.RemoveContainer(*deployment.ContainerId); removeErr != nil {
		return consumedPayload.DeploymentId, deployment.LastDeploymentInitiatedAt, removeErr
	}

	fmt.Println("container removed ", deployment.ContainerId)
	if _, updateErr := worker.deploymentService.UpdateDeployment(consumedPayload.DeploymentId, map[string]interface{}{
		"container_id": nil,
	}, context.Background()); updateErr != nil {
		return consumedPayload.DeploymentId, deployment.LastDeploymentInitiatedAt, updateErr
	}

	containerId, runErr := worker.dockerService.RunContainer(
		*deployment.DockerImageTag,
		deployment.Env,
		*port,
	)
	if runErr != nil {
		return consumedPayload.DeploymentId, deployment.LastDeploymentInitiatedAt, runErr
	}
	fmt.Println("docker image run successfully...")

	_, updateErr := worker.deploymentService.UpdateDeployment(
		consumedPayload.DeploymentId,
		map[string]interface{}{
			"latest_status":    enums.LIVE,
			"last_deployed_at": time.Now(),
			"container_id":     containerId,
		},
		context.Background(),
	)

	if updateErr != nil {
		return consumedPayload.DeploymentId, nil, updateErr
	}
	return consumedPayload.DeploymentId, deployment.LastDeploymentInitiatedAt, nil
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
