package payloads

type RunRepoWorkerPayload struct {
	DeploymentId   string                  `json:"deployment_id"`
	DockerImageTag string                  `bson:"docker_image_tag"`
	Env            *map[string]interface{} `json:"env"`
}
