package controller

import "github.com/gin-gonic/gin"

type AuthController struct {
}

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (controller *AuthController) SignUp(context *gin.Context) {

}

func (controller *AuthController) Login(context *gin.Context) {

}
