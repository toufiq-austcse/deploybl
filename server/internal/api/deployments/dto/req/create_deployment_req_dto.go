package req

import (
	"github.com/gin-gonic/gin"
)

type CreateDeploymentReqDto struct {
	Title          string             `json:"title" binding:"required"`
	RepositoryUrl  string             `json:"repository_url" binding:"required"`
	BranchName     string             `json:"branch_name" binding:"required"`
	RootDir        *string            `json:"root_dir"`
	DockerFilePath *string            `json:"docker_file_path"`
	Env            *map[string]string `json:"env"`
}

func (model *CreateDeploymentReqDto) Validate(c *gin.Context) error {
	err := c.BindJSON(model)
	if err != nil {
		return err
	}
	return nil
}
