package service

import (
	"fmt"
	"strings"

	"github.com/toufiq-austcse/deployit/pkg/cmd_runner"

	"github.com/toufiq-austcse/deployit/pkg/app_errors"

	deployItConfig "github.com/toufiq-austcse/deployit/config"
)

type DockerService struct{}

func NewDockerService() *DockerService {
	return &DockerService{}
}

func (dockerService *DockerService) RemoveContainer(containerId string) error {
	_, err := cmd_runner.RunCommand("docker", []string{"rm", "-f", containerId[:5]})
	if err != nil {
		return err
	}
	return nil
}

func (dockerService *DockerService) RunContainer(
	imageTag string,
	env *map[string]string,
	port *string,
) (*string, error) {
	args := []string{"run", "--network", deployItConfig.AppConfig.TRAEFIK_NETWORK_NAME}
	for k, v := range *env {
		args = append(args, "-e", fmt.Sprintf("%s=%s", k, v))
	}

	if port != nil {
		args = append(
			args,
			"-e",
			fmt.Sprintf("PORT=%s", *port),
			"-p",
			fmt.Sprintf(":%s", *port),
		)
	}
	args = append(args, "-d", imageTag)

	output, err := cmd_runner.RunCommand("docker", args)
	if err != nil {
		return nil, err
	}
	containerId := strings.Replace(*output, "\n", "", -1)
	fmt.Printf("Container ID: %s\n", containerId)
	return &containerId, nil
}

func (dockerService *DockerService) ListRunningContainerIds() ([]string, error) {
	containerIds := []string{}

	output, err := cmd_runner.RunCommand(
		"docker",
		[]string{"ps", "--no-trunc", "--format", "{{.ID}}"},
	)
	if err != nil {
		return nil, err
	}

	containerIds = strings.Split(*output, "\n")
	if len(containerIds) > 0 {
		containerIds = containerIds[:len(containerIds)-1]
	}

	return containerIds, nil
}

func (dockerService *DockerService) StopContainer(containerId string) error {
	_, err := cmd_runner.RunCommand("docker", []string{"stop", containerId[:5]})
	if err != nil {
		return err
	}
	return nil
}

func (dockerService *DockerService) GetTcpPort(containerID string) (*string, error) {
	output, err := cmd_runner.RunCommand(
		"docker",
		[]string{"exec", containerID, "netstat", "-tulpn"},
	)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(*output, "\n")
	if len(lines) < 3 {
		return nil, app_errors.ContainerPortNotFoundError
	}
	for i := 2; i < len(lines); i++ {
		if !strings.Contains(lines[i], "tcp") {
			continue
		}
		fields := strings.Fields(lines[i])
		if len(fields) < 4 {
			continue
		}
		localAddress := fields[3]
		if strings.Contains(localAddress, ":::") || strings.Contains(localAddress, "0.0.0.0") {
			parts := strings.Split(localAddress, ":")
			if len(parts) == 0 {
				continue
			}
			return &parts[len(parts)-1], nil
		}
	}

	return nil, app_errors.ContainerPortNotFoundError
}

func (dockerService *DockerService) BuildImage(
	dockerFilePath string,
	localDir string,
	dockerImageTag string,
	labels map[string]string,
) (*string, error) {
	args := []string{
		"build", "-f", dockerFilePath, localDir, "-t", dockerImageTag,
	}
	for k, v := range labels {
		args = append(args, "--label", fmt.Sprintf("%s=%s", k, v))
	}
	output, err := cmd_runner.RunCommand("docker", args)
	if err != nil {
		return nil, err
	}
	return output, nil
}
