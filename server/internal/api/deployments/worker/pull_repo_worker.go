package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	deployItConfig "github.com/toufiq-austcse/deployit/config"
	"github.com/toufiq-austcse/deployit/enums"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/mapper"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/model"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/service"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/worker/payloads"
	"github.com/toufiq-austcse/deployit/pkg/rabbit_mq"
	"github.com/toufiq-austcse/deployit/pkg/utils"
)

type PullRepoWorker struct {
	config            amqp.Config
	deploymentService *service.DeploymentService
	buildRepoWorker   *BuildRepoWorker
}

func NewPullRepoWorker(
	deploymentService *service.DeploymentService,
	buildRepoWorker *BuildRepoWorker,
) *PullRepoWorker {
	return &PullRepoWorker{
		config: rabbit_mq.New(deployItConfig.AppConfig.RABBIT_MQ_CONFIG.EXCHANGE,
			"topic",
			deployItConfig.AppConfig.RABBIT_MQ_CONFIG.REPOSITORY_PULL_QUEUE,
			deployItConfig.AppConfig.RABBIT_MQ_CONFIG.REPOSITORY_PULL_ROUTING_KEY),
		deploymentService: deploymentService,
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
		deploymentId, err := worker.ProcessPullRepoMessage(msg)
		if err != nil {
			fmt.Println("error in processing pull repo message ", err.Error())

			if deploymentId != "" {
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

func (worker *PullRepoWorker) ProcessPullRepoMessage(msg *message.Message) (string, error) {
	defer msg.Ack()

	consumedPayload := payloads.PullRepoWorkerPayload{}
	if err := json.Unmarshal(msg.Payload, &consumedPayload); err != nil {
		return "", err
	}

	if _, updateErr := worker.deploymentService.UpdateDeployment(consumedPayload.DeploymentId, map[string]interface{}{
		"latest_status": enums.PULLING,
	}, context.Background()); updateErr != nil {
		return consumedPayload.DeploymentId, updateErr
	}

	localRepoDir := utils.GetLocalRepoPath(consumedPayload.DeploymentId, consumedPayload.BranchName)
	if removeErr := os.RemoveAll(localRepoDir); removeErr != nil {
		return consumedPayload.DeploymentId, removeErr
	}

	if cloneError := worker.CloneRepo(consumedPayload.GitUrl, consumedPayload.BranchName, localRepoDir); cloneError != nil {
		return consumedPayload.DeploymentId, cloneError
	}
	fmt.Println("repository cloned successfully...")

	if buildRepoWorkPublishErr := worker.PublishBuildRepoWork(consumedPayload); buildRepoWorkPublishErr != nil {
		return consumedPayload.DeploymentId, buildRepoWorkPublishErr
	}

	if _, updateErr := worker.deploymentService.UpdateLatestStatus(consumedPayload.DeploymentId, enums.BUILDING, context.Background()); updateErr != nil {
		return consumedPayload.DeploymentId, updateErr
	}

	return consumedPayload.DeploymentId, nil
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

func (worker *PullRepoWorker) CloneRepo(gitUrl, branch, path string) error {
	cmd := exec.Command("git", "clone", "-b", branch, gitUrl, path)
	var out bytes.Buffer
	cmd.Stdout = &out

	fmt.Println("executing " + cmd.String())
	return cmd.Run()
}

func (worker *PullRepoWorker) PublishPullRepoWork(deployment *model.Deployment) {
	pullRepoWorkerPayload := mapper.ToPullRepoWorkerPayload(deployment)
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
