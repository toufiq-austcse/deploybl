package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/toufiq-austcse/go-api-boilerplate/config"
	"net/http"
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
		context.JSON(http.StatusOK, gin.H{
			"message": config.AppConfig.APP_NAME + " is Running",
		})
	}
}
