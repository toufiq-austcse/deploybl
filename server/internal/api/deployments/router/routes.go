package router

import (
	"github.com/gin-gonic/gin"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/controller"
)

func Setup(group *gin.RouterGroup, controller *controller.DeploymentController) {
	group.GET("", controller.DeploymentIndex)
	group.POST("", controller.DeploymentCreate)
	group.PATCH(":id", controller.DeploymentUpdate)
	group.PATCH(":id/env", controller.EnvUpdate)
	group.GET(":id", controller.DeploymentShow)
	group.POST(":id/restart", controller.DeploymentRestart)
	group.POST(":id/rebuild-and-redeploy", controller.DeploymentRebuildAndReDeploy)
	group.POST(":id/stop", controller.DeploymentStop)
	group.GET("latest-status", controller.DeploymentLatestStatus)
}

func SetupDeploymentCronRouter(group *gin.RouterGroup, controller *controller.DeploymentController) {
	group.POST("check-stopped-cron", controller.LiveCheckCron)

}
