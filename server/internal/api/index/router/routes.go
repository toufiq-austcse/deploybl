package router

import (
	"github.com/gin-gonic/gin"
	"github.com/toufiq-austcse/deployit/internal/api/index/controller"
)

func Setup(group *gin.RouterGroup, controller *controller.HealthController) {
	group.GET("", controller.Index)
}
