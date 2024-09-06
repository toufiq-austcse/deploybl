package req

import "github.com/gin-gonic/gin"

type CreateDeploymentReqDto struct {
	RepositoryProvider string                  `json:"repository_provider" binding:"required,oneof=github"`
	RepositoryUrl      string                  `json:"repository_url" binding:"required"`
	BranchName         string                  `json:"branch_name" binding:"required"`
	RootDir            *string                 `json:"root_dir"`
	DockerFilePath     *string                 `json:"docker_file_path"`
	Env                *map[string]interface{} `json:"env"`
}

func (model *CreateDeploymentReqDto) Validate(c *gin.Context) error {
	err := c.BindJSON(model)
	if err != nil {
		return err
	}
	return nil
}
