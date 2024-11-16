package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	deployItConfig "github.com/toufiq-austcse/deployit/config"
	"github.com/toufiq-austcse/deployit/enums"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/mapper"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/model"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/service"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/worker/payloads"
	"github.com/toufiq-austcse/deployit/pkg/cmd_runner"
	"github.com/toufiq-austcse/deployit/pkg/rabbit_mq"
	"github.com/toufiq-austcse/deployit/pkg/utils"
)

type PullRepoWorker struct {
	config            amqp.Config
	deploymentService *service.DeploymentService
	eventService      *service.EventService
	buildRepoWorker   *BuildRepoWorker
}

func NewPullRepoWorker(
	deploymentService *service.DeploymentService,
	eventService *service.EventService,
	buildRepoWorker *BuildRepoWorker,
) *PullRepoWorker {
	return &PullRepoWorker{
		config: rabbit_mq.New(deployItConfig.AppConfig.RABBIT_MQ_CONFIG.EXCHANGE,
			"topic",
			deployItConfig.AppConfig.RABBIT_MQ_CONFIG.REPOSITORY_PULL_QUEUE,
			deployItConfig.AppConfig.RABBIT_MQ_CONFIG.REPOSITORY_PULL_ROUTING_KEY),
		deploymentService: deploymentService,
		eventService:      eventService,
		buildRepoWorker:   buildRepoWorker,
	}
}

func (worker *PullRepoWorker) InitPullRepoSubscriber() {
	subscriber, err := amqp.NewSubscriber(worker.config, watermill.NewStdLogger(false, false))
	if err != nil {
		fmt.Println("error in pull repo subscriber ", err.Error())
		return
	}
	fmt.Println("PullRepoWorker Initialized")
	messages, err := subscriber.Subscribe(
		context.Background(),
		deployItConfig.AppConfig.RABBIT_MQ_CONFIG.REPOSITORY_PULL_QUEUE,
	)
	if err != nil {
		fmt.Println("error in pull repo subscriber ", err.Error())
		return
	}
	go worker.ProcessPullRepoMessages(messages)
}

func (worker *PullRepoWorker) ProcessPullRepoMessages(messages <-chan *message.Message) {
	for msg := range messages {
		deploymentId, event, logs, err := worker.ProcessPullRepoMessage(msg)
		if err != nil {
			fmt.Println("error in processing pull repo message ", err.Error())
			if logs != nil {
				worker.eventService.WriteEventLogToFile(*logs, event)
			}
			worker.eventService.WriteEventLogToFile("error in pulling repository", event)

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

func (worker *PullRepoWorker) ProcessPullRepoMessage(msg *message.Message) (string, *model.Event, *string, error) {
	defer msg.Ack()

	consumedPayload := payloads.PullRepoWorkerPayload{}
	if err := json.Unmarshal(msg.Payload, &consumedPayload); err != nil {
		return "", nil, nil, err
	}
	existingEvent, err := worker.eventService.FindById(consumedPayload.EventId)
	if err != nil {
		fmt.Println("error in finding event ", err.Error())
	}

	if _, updateErr := worker.deploymentService.UpdateLatestStatus(consumedPayload.DeploymentId, enums.PULLING, existingEvent, context.Background()); updateErr != nil {
		return consumedPayload.DeploymentId, nil, nil, updateErr
	}
	worker.eventService.WriteEventLogToFile(
		"cloning repository from "+consumedPayload.GitUrl+" branch "+consumedPayload.BranchName,
		existingEvent,
	)

	localRepoDir := utils.GetLocalRepoPath(consumedPayload.DeploymentId, consumedPayload.BranchName)
	if removeErr := os.RemoveAll(localRepoDir); removeErr != nil {
		return consumedPayload.DeploymentId, existingEvent, nil, removeErr
	}

	logs, cloneError := worker.CloneRepo(consumedPayload.GitUrl, consumedPayload.BranchName, localRepoDir)
	if cloneError != nil {
		return consumedPayload.DeploymentId, existingEvent, logs, cloneError
	}
	fmt.Println("repository cloned successfully...")
	worker.eventService.WriteEventLogToFile("cloned successfully", existingEvent)

	if buildRepoWorkPublishErr := worker.PublishBuildRepoWork(consumedPayload); buildRepoWorkPublishErr != nil {
		return consumedPayload.DeploymentId, existingEvent, logs, buildRepoWorkPublishErr
	}

	if _, updateErr := worker.deploymentService.UpdateLatestStatus(consumedPayload.DeploymentId, enums.BUILDING, existingEvent, context.Background()); updateErr != nil {
		return consumedPayload.DeploymentId, existingEvent, logs, updateErr
	}

	return consumedPayload.DeploymentId, existingEvent, logs, nil
}

func (worker *PullRepoWorker) PublishPullRepoJob(
	workerPayload payloads.PullRepoWorkerPayload,
) error {
	publisher, err := amqp.NewPublisher(worker.config, watermill.NewStdLogger(false, false))
	if err != nil {
		return err
	}

	payload, _ := json.Marshal(workerPayload)
	msg := message.NewMessage(watermill.NewUUID(), payload)
	err = publisher.Publish(deployItConfig.AppConfig.RABBIT_MQ_CONFIG.REPOSITORY_PULL_QUEUE, msg)
	if err != nil {
		return err
	}
	return publisher.Close()
}

func (worker *PullRepoWorker) CloneRepo(gitUrl, branch, path string) (*string, error) {
	logs, err := cmd_runner.RunCommand("git", []string{"clone", "-b", branch, gitUrl, path})
	if err != nil {
		return logs, err
	}
	return logs, nil
}

func (worker *PullRepoWorker) PublishPullRepoWork(
	deployment *model.Deployment,
	event *model.Event,
) {
	pullRepoWorkerPayload := mapper.ToPullRepoWorkerPayload(deployment, event)
	fmt.Println("publishing pull repo worker payload ", pullRepoWorkerPayload)
	err := worker.PublishPullRepoJob(pullRepoWorkerPayload)
	if err != nil {
		fmt.Println("error while publishing pull repo worker job ", err.Error())
		return
	}
}

func (worker *PullRepoWorker) PublishBuildRepoWork(
	pullRepoJob payloads.PullRepoWorkerPayload,
) error {
	buildRepoWorkerPayload := mapper.ToBuildRepoWorkerPayload(pullRepoJob)
	fmt.Println("Publishing ", buildRepoWorkerPayload)
	return worker.buildRepoWorker.PublishBuildRepoJob(buildRepoWorkerPayload)
}
