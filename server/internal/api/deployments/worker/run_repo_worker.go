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
	"github.com/toufiq-austcse/deployit/pkg/rabbit_mq"
	"os/exec"
)

type RunRepoWorker struct {
	config            amqp.Config
	deploymentService *service.DeploymentService
}

func NewRunRepoWorker(deploymentService *service.DeploymentService) *RunRepoWorker {
	return &RunRepoWorker{config: rabbit_mq.New(deployItConfig.AppConfig.RABBIT_MQ_CONFIG.EXCHANGE,
		"topic",
		deployItConfig.AppConfig.RABBIT_MQ_CONFIG.REPOSITORY_RUN_QUEUE,
		deployItConfig.AppConfig.RABBIT_MQ_CONFIG.REPOSITORY_RUN_ROUTING_KEY),
		deploymentService: deploymentService}
}

func (worker *RunRepoWorker) InitRunRepoSubscriber() {
	subscriber, err := amqp.NewSubscriber(worker.config, watermill.NewStdLogger(false, false))
	if err != nil {
		fmt.Println("error in run repo subscriber ", err.Error())
		return
	}
	fmt.Println("RunRepoSubscriber Initialized")
	messages, err := subscriber.Subscribe(context.Background(), deployItConfig.AppConfig.RABBIT_MQ_CONFIG.REPOSITORY_PULL_QUEUE)
	if err != nil {
		fmt.Println("error in run repo subscriber ", err.Error())
		return
	}
	go worker.ProcessRunRepoMessage(messages)

}
func (worker *RunRepoWorker) ProcessRunRepoMessage(messages <-chan *message.Message) {
	consumedPayload := payloads.RunRepoWorkerPayload{}
	for msg := range messages {
		err := json.Unmarshal(msg.Payload, &consumedPayload)
		if err != nil {
			fmt.Println("error in parsing rabbitmq message in pull repo worker", err)
			msg.Ack()
			continue
		}
		fmt.Println("consumed run job ", consumedPayload)
		runErr := worker.RunRepo(consumedPayload)
		if runErr != nil {
			_, updateErr := worker.deploymentService.UpdateDeployment(consumedPayload.DeploymentId, map[string]interface{}{
				"latest_status": enums.FAILED,
			}, context.Background())
			if updateErr != nil {
				fmt.Println("error while updating status... ", updateErr.Error())
			}

			msg.Ack()
			continue
		}
		fmt.Println("docker image run successfully...")

		// we need to Acknowledge that we received and processed the message,
		// otherwise, it will be resent over and over again.
		msg.Ack()
	}
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

func (worker *RunRepoWorker) RunRepo(payload payloads.RunRepoWorkerPayload) error {
	//envString := worker.GetEnvString(*payload.Env)
	imageID := "66e1ecce5e1128018ac9bd88"
	port := "4000"

	// Construct the command
	cmd := exec.Command(
		"/usr/local/bin/docker", "run",
		"-e", fmt.Sprintf("PORT=%s", port), // No need for single quotes around env variables in exec.Command
		"-p", fmt.Sprintf("%s:%s", port, port),
		"-d", imageID,
	)

	// Run the command and capture the output
	output, err := cmd.CombinedOutput()

	// Handle errors if the command fails
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		fmt.Printf("Output: %s\n", string(output))
		return err
	}

	// Print the container ID (which is the output of the docker run -d command)
	fmt.Printf("Container ID: %s\n", string(output))
	return nil
}
func (worker *RunRepoWorker) GetEnvString(env map[string]interface{}) string {
	envString := "'PORT=4000'"
	//for k, v := range env {
	//	envString += " -e " + k + "=" + v.(string)
	//}
	return envString
}
