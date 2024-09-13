package payloads

type PullRepoWorkerPayload struct {
	DeploymentId   string             `json:"deployment_id"`
	BranchName     string             `json:"branch_name"`
	GitUrl         string             `json:"git_url"`
	DockerFilePath string             `json:"docker_file_path"`
	Env            *map[string]string `json:"env"`
}
