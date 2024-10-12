package router

import (
	"github.com/gin-gonic/gin"
	"github.com/toufiq-austcse/deployit/internal/api/repositories/controller"
)

func Setup(group *gin.RouterGroup, controller *controller.RepoController) {
	group.GET("", controller.GetRepoDetails)
	group.GET("branches", controller.GetRepoBranches)
}
