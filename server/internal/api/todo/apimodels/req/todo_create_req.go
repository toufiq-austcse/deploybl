package req

import "github.com/gin-gonic/gin"

type TodoCreateReqModel struct {
	Title  string `json:"title" binding:"required"`
	Status string `json:"status" binding:"required"`
}

func (model *TodoCreateReqModel) Validate(c *gin.Context) error {
	err := c.BindJSON(model)
	if err != nil {
		return err
	}
	return nil
}
