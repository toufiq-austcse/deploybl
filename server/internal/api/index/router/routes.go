package router

import (
	"github.com/gin-gonic/gin"
	"github.com/toufiq-austcse/go-api-boilerplate/internal/api/index/controller"
)

func Setup(group *gin.RouterGroup) {
	group.GET("", controller.Index())

}
