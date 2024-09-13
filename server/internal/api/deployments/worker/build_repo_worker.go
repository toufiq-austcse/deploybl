package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	"os/exec"
)

type BuildRepoWorker struct {
	config            amqp.Config
	deploymentService *service.DeploymentService
	runRepoWorker     *RunRepoWorker
}

func NewBuildRepoWorker(deploymentService *service.DeploymentService, runRepoWorker *RunRepoWorker) *BuildRepoWorker {
	return &BuildRepoWorker{config: rabbit_mq.New(deployItConfig.AppConfig.RABBIT_MQ_CONFIG.EXCHANGE,
		"topic",
		deployItConfig.AppConfig.RABBIT_MQ_CONFIG.REPOSITORY_BUILD_QUEUE,
		deployItConfig.AppConfig.RABBIT_MQ_CONFIG.REPOSITORY_BUILD_ROUTING_KEY),
		deploymentService: deploymentService,
		runRepoWorker:     runRepoWorker}
}

func (worker *BuildRepoWorker) InitBuildRepoSubscriber() {
	subscriber, err := amqp.NewSubscriber(worker.config, watermill.NewStdLogger(false, false))
	if err != nil {
		fmt.Println("error in build repo subscriber ", err.Error())
		return
	}
	fmt.Println("BuildRepoWorkerPayload Initialized")
	messages, err := subscriber.Subscribe(context.Background(), deployItConfig.AppConfig.RABBIT_MQ_CONFIG.REPOSITORY_PULL_QUEUE)
	if err != nil {
		fmt.Println("error in build repo subscriber ", err.Error())
		return
	}
	go worker.ProcessBuildRepoMessage(messages)

}
func (worker *BuildRepoWorker) ProcessBuildRepoMessage(messages <-chan *message.Message) {
	consumedPayload := payloads.BuildRepoWorkerPayload{}
	for msg := range messages {
		err := json.Unmarshal(msg.Payload, &consumedPayload)
		if err != nil {
			fmt.Println("error in parsing rabbitmq message in pull repo worker", err)
			msg.Ack()
			continue
		}
		fmt.Println("consumed build job ", consumedPayload)

		dockerImageTag, buildRepoErr := worker.BuildRepo(consumedPayload)
		if buildRepoErr != nil {
			_, updateErr := worker.deploymentService.UpdateDeployment(consumedPayload.DeploymentId, map[string]interface{}{
				"latest_status": enums.FAILED,
			}, context.Background())
			if updateErr != nil {
				fmt.Println("error while updating status... ", updateErr.Error())
			}

			msg.Ack()
			continue
		}
		fmt.Println("Docker image built successfully")
		_, updateErr := worker.deploymentService.UpdateDeployment(consumedPayload.DeploymentId, map[string]interface{}{
			"docker_image_tag": dockerImageTag,
		}, context.Background())

		if updateErr != nil {
			fmt.Println("error while updating image tag... ", updateErr.Error())
			msg.Ack()
			continue
		}

		publishRunRepoError := worker.PublishRunRepoWork(consumedPayload, *dockerImageTag)
		if publishRunRepoError != nil {
			fmt.Println("error in publishing run repo work ", publishRunRepoError.Error())
		}

		// we need to Acknowledge that we received and processed the message,
		// otherwise, it will be resent over and over again.
		msg.Ack()
	}
}

func (worker *BuildRepoWorker) PublishBuildRepoJob(workerPayload payloads.BuildRepoWorkerPayload) error {
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
	localDir := utils.GetLocalRepoPath(payload.DeploymentId)
	dockerFilePath := localDir
	if payload.DockerFilePath != "." {
		dockerFilePath += "/" + payload.DockerFilePath
	}
	dockerImageTag := payload.DeploymentId
	host := payload.SubDomainName + "." + deployItConfig.AppConfig.BASE_DOMAIN
	hostRule := fmt.Sprintf("HOST(`%s`)", host)
	labels := map[string]string{
		"traefik.enable":                 "true",
		"traefik.http.routers.test.rule": hostRule,
	}
	args := []string{
		"build", dockerFilePath, "-t", dockerImageTag,
	}
	for k, v := range labels {
		args = append(args, "--label", fmt.Sprintf("%s=%s", k, v))
	}
	cmd := exec.Command("docker", args...)
	var out bytes.Buffer
	cmd.Stdout = &out

	fmt.Println("executing " + cmd.String())
	err := cmd.Run()
	if err != nil {
		fmt.Println("docker build err", err.Error())
		return nil, err
	}
	fmt.Println(out.String())
	return &dockerImageTag, nil
}

func (worker *BuildRepoWorker) PublishRunRepoWork(buildRepoWorkerPayload payloads.BuildRepoWorkerPayload, dockerImageTag string) error {
	runRepoWorkerPayload := mapper.ToRunRepoWorkerPayload(buildRepoWorkerPayload, dockerImageTag)
	fmt.Println("Publishing ", runRepoWorkerPayload)
	return worker.runRepoWorker.PublishRunRepoJob(runRepoWorkerPayload)
}
