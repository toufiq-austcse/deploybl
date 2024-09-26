package payloads

type PullRepoWorkerPayload struct {
	DeploymentId   string             `json:"deployment_id"`
	BranchName     string             `json:"branch_name"`
	SubDomainName  string             `json:"sub_domain_name"`
	GitUrl         string             `json:"git_url"`
	RootDir        *string            `json:"root_dir"`
	DockerFilePath string             `json:"docker_file_path"`
	Env            *map[string]string `json:"env"`
}
