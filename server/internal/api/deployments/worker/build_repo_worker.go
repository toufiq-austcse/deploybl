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
	"github.com/toufiq-austcse/deployit/internal/api/deployments/mapper"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/service"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/worker/payloads"
	"github.com/toufiq-austcse/deployit/pkg/rabbit_mq"
	"github.com/toufiq-austcse/deployit/pkg/utils"
)

type BuildRepoWorker struct {
	config            amqp.Config
	deploymentService *service.DeploymentService
	preRunRepoWorker  *PreRunRepoWorker
	dockerService     *service.DockerService
	eventService      *service.EventService
}

func NewBuildRepoWorker(
	deploymentService *service.DeploymentService,
	preRunRepoWorker *PreRunRepoWorker,
	dockerService *service.DockerService,
	eventService *service.EventService,
) *BuildRepoWorker {
	return &BuildRepoWorker{
		config: rabbit_mq.New(deployItConfig.AppConfig.RABBIT_MQ_CONFIG.EXCHANGE,
			"topic",
			deployItConfig.AppConfig.RABBIT_MQ_CONFIG.REPOSITORY_BUILD_QUEUE,
			deployItConfig.AppConfig.RABBIT_MQ_CONFIG.REPOSITORY_BUILD_ROUTING_KEY),
		deploymentService: deploymentService,
		preRunRepoWorker:  preRunRepoWorker,
		dockerService:     dockerService,
		eventService:      eventService,
	}
}

func (worker *BuildRepoWorker) InitBuildRepoSubscriber() {
	subscriber, err := amqp.NewSubscriber(worker.config, watermill.NewStdLogger(false, false))
	if err != nil {
		fmt.Println("error in build repo subscriber ", err.Error())
		return
	}
	fmt.Println("BuildRepoWorkerPayload Initialized")
	messages, err := subscriber.Subscribe(
		context.Background(),
		deployItConfig.AppConfig.RABBIT_MQ_CONFIG.REPOSITORY_PULL_QUEUE,
	)
	if err != nil {
		fmt.Println("error in build repo subscriber ", err.Error())
		return
	}
	go worker.ProcessBuildRepoMessages(messages)
}

func (worker *BuildRepoWorker) ProcessBuildRepoMessages(messages <-chan *message.Message) {
	for msg := range messages {
		deploymentId, event, buildLogs, err := worker.ProcessBuildRepoMessage(msg)
		fmt.Println("buildLogs ", buildLogs)
		if err != nil {
			fmt.Println("error in processing build repo message ", err.Error())
			if buildLogs != nil {
				worker.eventService.WriteEventLogToFile(*buildLogs, event)
			}
			worker.eventService.WriteEventLogToFile("error in building docker image: "+err.Error(), event)

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

func (worker *BuildRepoWorker) ProcessBuildRepoMessage(msg *message.Message) (string, *model.Event, *string, error) {
	defer msg.Ack()

	consumedPayload := payloads.BuildRepoWorkerPayload{}
	err := json.Unmarshal(msg.Payload, &consumedPayload)
	if err != nil {
		return "", nil, nil, err
	}
	fmt.Println("consumed build job ", consumedPayload)
	event, err := worker.eventService.FindById(consumedPayload.EventId)
	if err != nil {
		fmt.Println("error in finding event ", err.Error())
	}

	worker.eventService.WriteEventLogToFile("building dockerfile", event)
	dockerImageTag, buildLogs, buildRepoErr := worker.BuildRepo(consumedPayload)
	if buildRepoErr != nil {
		return consumedPayload.DeploymentId, event, buildLogs, buildRepoErr
	}
	fmt.Println("Docker image built successfully")
	worker.eventService.WriteEventLogToFile(*buildLogs, event)
	worker.eventService.WriteEventLogToFile("docker image built successfully", event)
	if _, updateErr := worker.deploymentService.UpdateDeployment(consumedPayload.DeploymentId, map[string]interface{}{
		"docker_image_tag": dockerImageTag,
	}, event, context.Background()); updateErr != nil {
		return consumedPayload.DeploymentId, event, buildLogs, updateErr
	}

	if publishRunRepoError := worker.PublishPreRunRepoWork(consumedPayload); publishRunRepoError != nil {
		return consumedPayload.DeploymentId, event, buildLogs, publishRunRepoError
	}

	return consumedPayload.DeploymentId, event, buildLogs, nil
}

func (worker *BuildRepoWorker) PublishBuildRepoJob(
	workerPayload payloads.BuildRepoWorkerPayload,
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

func (worker *BuildRepoWorker) BuildRepo(payload payloads.BuildRepoWorkerPayload) (*string, *string, error) {
	localDir := utils.GetLocalRepoPath(payload.DeploymentId, payload.BranchName)
	dockerFilePath := localDir
	if payload.DockerFilePath != "." {
		dockerFilePath += "/" + payload.DockerFilePath
	}
	dockerImageTag := payload.DeploymentId
	host := payload.SubDomainName + "." + deployItConfig.AppConfig.BASE_DOMAIN
	hostRule := fmt.Sprintf("HOST(`%s`)", host)
	labels := map[string]string{
		"traefik.enable": "true",
		fmt.Sprintf("traefik.http.routers.%s.rule", payload.DeploymentId): hostRule,
	}
	buildLogs, err := worker.dockerService.BuildImage(dockerFilePath, localDir, dockerImageTag, labels)
	if err != nil {
		return nil, buildLogs, err
	}

	return &dockerImageTag, buildLogs, nil
}

func (worker *BuildRepoWorker) PublishPreRunRepoWork(
	buildRepoWorkerPayload payloads.BuildRepoWorkerPayload,
) error {
	preRunRepoWorkerPayload := mapper.ToPreRunRepoWorkerPayload(buildRepoWorkerPayload)
	fmt.Println("Publishing ", preRunRepoWorkerPayload)
	return worker.preRunRepoWorker.PublishPreRunRepoJob(preRunRepoWorkerPayload)
}
