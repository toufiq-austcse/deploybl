package req

import "github.com/gin-gonic/gin"

type SignUpReqDto struct {
	Name     string `json:"name"     binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,gt=5"`
}

func (model *SignUpReqDto) Validate(c *gin.Context) error {
	err := c.BindJSON(model)
	if err != nil {
		return err
	}
	return nil
}
