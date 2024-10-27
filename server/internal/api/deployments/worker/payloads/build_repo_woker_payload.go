package payloads

import "go.mongodb.org/mongo-driver/bson/primitive"

type BuildRepoWorkerPayload struct {
	DeploymentId   string             `json:"deployment_id"`
	SubDomainName  string             `json:"sub_domain_name"`
	BranchName     string             `json:"branch_name"`
	RootDir        *string            `json:"root_dir"`
	DockerFilePath string             `json:"docker_file_path"`
	Env            *map[string]string `json:"env"`
	EventId        primitive.ObjectID `json:"event_id"`
}
