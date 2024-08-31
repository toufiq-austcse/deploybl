package res

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type DeploymentRes struct {
	Id                 primitive.ObjectID `json:"_id"`
	Title              string             `json:"title"`
	LatestStatus       string             `json:"latest_status"`
	LastDeployedAt     *time.Time         `json:"last_deployed_at"`
	RepositoryProvider string             `json:"repository_provider"`
	BranchName         string             `json:"branch_name"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
}

type DeploymentDetailsRes struct {
	Id                 primitive.ObjectID      `json:"_id"`
	Title              string                  `json:"title"`
	SubDomainName      string                  `json:"sub_domain_name"`
	LatestStatus       string                  `json:"latest_status"`
	LastDeployedAt     *time.Time              `json:"last_deployed_at"`
	RepositoryProvider string                  `json:"repository_provider"`
	RepositoryUrl      string                  `json:"repository_url"`
	BranchName         string                  `json:"branch_name"`
	DockerFilePath     string                  `json:"docker_file_path"`
	DockerImageTag     *string                 `json:"docker_image_tag"`
	ContainerId        *string                 `json:"container_id"`
	Env                *map[string]interface{} `json:"env"`
	CreatedAt          time.Time               `json:"created_at"`
	UpdatedAt          time.Time               `json:"updated_at"`
}
