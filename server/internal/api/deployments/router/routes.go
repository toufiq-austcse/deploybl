package router

import (
	"github.com/gin-gonic/gin"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/controller"
)

func Setup(group *gin.RouterGroup, controller *controller.DeploymentController) {
	group.GET("", controller.DeploymentIndex)
	group.POST("", controller.DeploymentCreate)
	group.PATCH(":id", controller.DeploymentUpdate)
	group.PUT(":id/env", controller.EnvUpdate)
	group.GET(":id", controller.DeploymentShow)
	group.GET("latest-status", controller.DeploymentLatestStatus)

}
