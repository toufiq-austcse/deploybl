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
	"github.com/toufiq-austcse/deployit/internal/api/deployments/model"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/service"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/worker/payloads"
	"github.com/toufiq-austcse/deployit/pkg/rabbit_mq"
	"os/exec"
)

type PullRepoWorker struct {
	config            amqp.Config
	deploymentService *service.DeploymentService
}

func NewPullRepoWorker(deploymentService *service.DeploymentService) *PullRepoWorker {
	return &PullRepoWorker{config: rabbit_mq.New(deployItConfig.AppConfig.RABBIT_MQ_CONFIG.EXCHANGE,
		"topic",
		deployItConfig.AppConfig.RABBIT_MQ_CONFIG.REPOSITORY_PULL_QUEUE,
		deployItConfig.AppConfig.RABBIT_MQ_CONFIG.REPOSITORY_PULL_ROUTING_KEY),
		deploymentService: deploymentService}
}

func (worker *PullRepoWorker) InitPullRepoSubscriber() {
	subscriber, err := amqp.NewSubscriber(worker.config, watermill.NewStdLogger(false, false))
	if err != nil {
		fmt.Println("error in pull repo subscriber ", err.Error())
		return
	}
	fmt.Println("PullRepoWorker Initialized")
	messages, err := subscriber.Subscribe(context.Background(), deployItConfig.AppConfig.RABBIT_MQ_CONFIG.REPOSITORY_PULL_QUEUE)
	if err != nil {
		fmt.Println("error in pull repo subscriber ", err.Error())
		return
	}
	go worker.ProcessMessage(messages)

}
func (worker *PullRepoWorker) ProcessMessage(messages <-chan *message.Message) {
	consumedPayload := payloads.PullRepoWorkerPayload{}
	for msg := range messages {
		err := json.Unmarshal(msg.Payload, &consumedPayload)
		if err != nil {
			fmt.Println("error in parsing rabbitmq message in pull repo worker", err)
			continue
		}
		_, updateErr := worker.deploymentService.UpdateStatus(consumedPayload.DeploymentId, enums.PULLING, context.Background())
		if updateErr != nil {
			fmt.Println("error while updating status... ", updateErr.Error())
			continue
		}
		cloneError := worker.CloneRepo(consumedPayload.GitUrl, deployItConfig.AppConfig.REPOSITORIES_PATH+"/"+consumedPayload.DeploymentId)
		if cloneError != nil {
			_, updateErr = worker.deploymentService.UpdateStatus(consumedPayload.DeploymentId, enums.BUILDING, context.Background())
			if updateErr != nil {
				fmt.Println("error while updating status... ", updateErr.Error())
				continue
			}
		}

		// we need to Acknowledge that we received and processed the message,
		// otherwise, it will be resent over and over again.
		msg.Ack()
	}
}

func (worker *PullRepoWorker) PublishPullRepoJob(workerPayload payloads.PullRepoWorkerPayload) error {
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

func (worker *PullRepoWorker) CloneRepo(gitUrl, path string) error {
	cmd := exec.Command("git", "clone", gitUrl, path)
	var out bytes.Buffer
	cmd.Stdout = &out

	fmt.Println("executing " + cmd.String())
	err := cmd.Run()
	if err != nil {
		fmt.Println("git clone error ", err.Error())
		return err
	}
	fmt.Println(out.String())
	return nil
}

func (worker *PullRepoWorker) PublishPullRepoWork(deployment *model.Deployment) {
	pullRepoWorkerPayload := mapper.ToPullRepoWorkerPayload(deployment)
	err := worker.PublishPullRepoJob(pullRepoWorkerPayload)
	if err != nil {
		fmt.Println("error while publishing pull repo worker job ", err.Error())
		return
	}
}
