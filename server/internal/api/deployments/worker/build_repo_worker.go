package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"

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
}

func NewBuildRepoWorker(
	deploymentService *service.DeploymentService,
	preRunRepoWorker *PreRunRepoWorker,
) *BuildRepoWorker {
	return &BuildRepoWorker{
		config: rabbit_mq.New(deployItConfig.AppConfig.RABBIT_MQ_CONFIG.EXCHANGE,
			"topic",
			deployItConfig.AppConfig.RABBIT_MQ_CONFIG.REPOSITORY_BUILD_QUEUE,
			deployItConfig.AppConfig.RABBIT_MQ_CONFIG.REPOSITORY_BUILD_ROUTING_KEY),
		deploymentService: deploymentService,
		preRunRepoWorker:  preRunRepoWorker,
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
		deploymentId, err := worker.ProcessBuildRepoMessage(msg)
		if err != nil {
			fmt.Println("error in processing build repo message ", err.Error())

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

func (worker *BuildRepoWorker) ProcessBuildRepoMessage(msg *message.Message) (string, error) {
	defer msg.Ack()

	consumedPayload := payloads.BuildRepoWorkerPayload{}
	err := json.Unmarshal(msg.Payload, &consumedPayload)
	if err != nil {
		return "", err
	}
	fmt.Println("consumed build job ", consumedPayload)

	dockerImageTag, buildRepoErr := worker.BuildRepo(consumedPayload)
	if buildRepoErr != nil {
		return consumedPayload.DeploymentId, buildRepoErr
	}
	fmt.Println("Docker image built successfully")
	if _, updateErr := worker.deploymentService.UpdateDeployment(consumedPayload.DeploymentId, map[string]interface{}{
		"docker_image_tag": dockerImageTag,
	}, context.Background()); updateErr != nil {
		return consumedPayload.DeploymentId, updateErr
	}

	if publishRunRepoError := worker.PublishPreRunRepoWork(consumedPayload); publishRunRepoError != nil {
		return consumedPayload.DeploymentId, publishRunRepoError
	}

	return consumedPayload.DeploymentId, nil
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

func (worker *BuildRepoWorker) BuildRepo(payload payloads.BuildRepoWorkerPayload) (*string, error) {
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
	args := []string{
		"build", "-f", dockerFilePath, localDir, "-t", dockerImageTag,
	}
	for k, v := range labels {
		args = append(args, "--label", fmt.Sprintf("%s=%s", k, v))
	}
	cmd := exec.Command("docker", args...)
	var out bytes.Buffer
	var stdErr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stdErr

	fmt.Println("executing " + cmd.String())
	err := cmd.Run()
	if err != nil {
		return nil, errors.New(stdErr.String())
	}
	fmt.Println(out.String())
	return &dockerImageTag, nil
}

func (worker *BuildRepoWorker) PublishPreRunRepoWork(
	buildRepoWorkerPayload payloads.BuildRepoWorkerPayload,
) error {
	preRunRepoWorkerPayload := mapper.ToPreRunRepoWorkerPayload(buildRepoWorkerPayload)
	fmt.Println("Publishing ", preRunRepoWorkerPayload)
	return worker.preRunRepoWorker.PublishPreRunRepoJob(preRunRepoWorkerPayload)
}
