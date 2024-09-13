package payloads

type BuildRepoWorkerPayload struct {
	DeploymentId   string             `json:"deployment_id"`
	SubDomainName  string             `json:"sub_domain_name"`
	DockerFilePath string             `json:"docker_file_path"`
	Env            *map[string]string `json:"env"`
}
