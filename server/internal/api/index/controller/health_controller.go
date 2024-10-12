package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/toufiq-austcse/deployit/config"
)

// Index hosts godoc
// @Summary  Health Check
// @Tags     Index
// @Accept   json
// @Produce  json
// @Success  200
// @Router   / [get]
func Index() gin.HandlerFunc {
	return func(context *gin.Context) {
		//err := worker.PublishPullRepoJob()
		//if err != nil {
		//	fmt.Println("error while publishing ", err)
		//}
		context.JSON(http.StatusOK, gin.H{
			"message": config.AppConfig.APP_NAME + " is Running",
		})
	}
}
