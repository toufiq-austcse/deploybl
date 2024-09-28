package req

import "github.com/gin-gonic/gin"

type UpdateDeploymentReqDto struct {
	Title          *string `json:"title"`
	BranchName     *string `json:"branch_name" `
	RootDir        *string `json:"root_dir"`
	DockerFilePath *string `json:"docker_file_path"`
}

func (model *UpdateDeploymentReqDto) Validate(c *gin.Context) error {
	err := c.BindJSON(model)
	if err != nil {
		return err
	}
	return nil
}
