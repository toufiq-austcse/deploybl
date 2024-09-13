package payloads

type BuildRepoWorkerPayload struct {
	DeploymentId   string             `json:"deployment_id"`
	DockerFilePath string             `json:"docker_file_path"`
	Env            *map[string]string `json:"env"`
}
