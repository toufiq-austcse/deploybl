package service

import (
	"bytes"
	"errors"
	"fmt"
	deployItConfig "github.com/toufiq-austcse/deployit/config"
	"os/exec"
	"strings"
)

type DockerService struct {
}

func NewDockerService() *DockerService {
	return &DockerService{}
}

func (dockerService *DockerService) RemoveContainer(containerId string) error {
	cmd := exec.Command("docker", "rm", "-f", containerId[:5])
	var out bytes.Buffer
	var err bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &err

	fmt.Println("executing " + cmd.String())
	runErr := cmd.Run()
	if runErr != nil {
		return errors.New(err.String())
	}
	return nil
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

func (dockerService *DockerService) ListRunningContainerIds() ([]string, error) {
	containerIds := []string{}

	cmd := exec.Command("docker", "ps", "--no-trunc", "--format", "{{.ID}}")
	var out bytes.Buffer
	cmd.Stdout = &out

	fmt.Println("executing " + cmd.String())
	err := cmd.Run()
	if err != nil {
		return []string{}, err
	}
	containerIdsStringRes := string(out.Bytes())
	containerIds = strings.Split(containerIdsStringRes, "\n")
	if len(containerIds) > 0 {
		containerIds = containerIds[:len(containerIds)-1]
	}

	return containerIds, nil
}
