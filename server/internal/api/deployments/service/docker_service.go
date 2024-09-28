package service

import (
	"bytes"
	"errors"
	"fmt"
	deployItConfig "github.com/toufiq-austcse/deployit/config"
	"os/exec"
)

type DockerService struct{}

func NewDockerService() *DockerService {
	return &DockerService{}
}

func (dockerService *DockerService) RemoveContainer(containerId string) (*string, error) {
	cmd := exec.Command("docker", "stop", containerId)
	var out bytes.Buffer
	var err bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &err

	fmt.Println("executing " + cmd.String())
	runErr := cmd.Run()
	if runErr != nil {
		return nil, errors.New(err.String())
	}
	output, combinedOutputErr := cmd.CombinedOutput()
	if combinedOutputErr != nil {
		fmt.Printf("Error: %s\n", err)
		fmt.Printf("Output: %s\n", string(output))
		return nil, combinedOutputErr
	}
	ouputContainerid := string(output)
	return &ouputContainerid, nil
}

func (dockerService *DockerService) RunContainer(imageTag string, env *map[string]string) (*string, error) {
	port := "4000"

	args := []string{"run", "--network", deployItConfig.AppConfig.TRAEFIK_NETWORK_NAME}
	for k, v := range *env {
		if k == "PORT" {
			continue
		}
		args = append(args, "-e", fmt.Sprintf("%s=%s", k, v))
	}

	args = append(args, "-e", fmt.Sprintf("PORT=%s", port), "-p", fmt.Sprintf(":%s", port), "-d", imageTag)

	// Construct the command
	cmd := exec.Command(
		"docker", args...,
	)
	fmt.Println("executing ", cmd.String())

	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("Error: %s\n", err)
		fmt.Printf("Output: %s\n", string(output))
		return nil, err
	}

	containerId := string(output)
	fmt.Printf("Container ID: %s\n", string(output))
	return &containerId, nil
}
