package payloads

type StopRepoWorkerPayload struct {
	DeploymentId string `json:"deployment_id"`
}
