package service

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/toufiq-austcse/deployit/pkg/app_errors"

	deployItConfig "github.com/toufiq-austcse/deployit/config"
)

type DockerService struct{}

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

func (dockerService *DockerService) RunContainer(
	imageTag string,
	env *map[string]string,
	port string,
) (*string, error) {
	args := []string{"run", "--network", deployItConfig.AppConfig.TRAEFIK_NETWORK_NAME}
	for k, v := range *env {
		if k == "PORT" {
			continue
		}
		args = append(args, "-e", fmt.Sprintf("%s=%s", k, v))
	}

	args = append(
		args,
		"-e",
		fmt.Sprintf("PORT=%s", port),
		"-p",
		fmt.Sprintf(":%s", port),
		"-d",
		imageTag,
	)

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

	outputString := string(output)
	containerId := strings.Replace(outputString, "\n", "", -1)
	fmt.Printf("Container ID: %s\n", string(output))
	return &containerId, nil
}

func (dockerService *DockerService) PreRun(
	imageTag string,
	env *map[string]string,
) (*string, error) {
	args := []string{"run"}
	for k, v := range *env {
		args = append(args, "-e", fmt.Sprintf("%s=%s", k, v))
	}

	args = append(
		args,
		"-d",
		imageTag,
	)

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

	outputString := string(output)
	containerId := strings.Replace(outputString, "\n", "", -1)

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

func (dockerService *DockerService) StopContainer(containerId string) error {
	cmd := exec.Command("docker", "stop", containerId[:5])
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

func (dockerService *DockerService) GetTcpPort(containerID string) (*string, error) {
	// Run docker exec with netstat command
	cmd := exec.Command(
		"docker",
		"exec",
		containerID,
		"netstat",
		"-tulpn",
	)

	var out bytes.Buffer
	cmd.Stdout = &out

	fmt.Println("executing " + cmd.String())
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	output := string(out.Bytes())
	lines := strings.Split(output, "\n")
	if len(lines) < 3 {
		return nil, nil
	}
	tcpLine := strings.Split(output, "\n")[2]
	fields := strings.Fields(tcpLine)
	if len(fields) < 4 {
		return nil, nil
	}
	localAddress := fields[3]
	parts := strings.Split(localAddress, ":")
	if len(parts) < 2 {
		return nil, app_errors.ContainerPortNotFoundError
	}
	return &parts[1], err
}
