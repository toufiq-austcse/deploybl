package router

import (
	"github.com/gin-gonic/gin"
	controller "github.com/toufiq-austcse/deployit/internal/api/auth/controller"
)

func Setup(
	group *gin.RouterGroup,
	controller *controller.AuthController,
) {
	group.POST("signup", controller.SignUp)
	group.POST("login", controller.Login)
}
