package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/toufiq-austcse/deployit/config"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/dto/req"
	"github.com/toufiq-austcse/deployit/pkg/api_response"
	"github.com/toufiq-austcse/deployit/pkg/http_clients/github"
	"net/http"
)

type DeploymentController struct {
	githubHttpClient *github.GithubHttpClient
}

func NewDeploymentController(githubHttpClient *github.GithubHttpClient) *DeploymentController {
	return &DeploymentController{
		githubHttpClient: githubHttpClient,
	}
}

// DeploymentIndex
// @Summary  Deployment Index
// @Tags     Deployments
// @Accept   json
// @Produce  json
// @Success  200
// @Router   /deployments [get]
func DeploymentIndex() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": config.AppConfig.APP_NAME + " is Running",
		})
	}
}

// DeploymentCreate
// @Summary  Create Deployment
// @Param    request  body      req.CreateDeploymentReqDto  true  "Create Deployment Body"
// @Tags     Deployments
// @Accept   json
// @Produce  json
// @Success  200
// @Router   /deployments [post]
func (controller *DeploymentController) DeploymentCreate(context *gin.Context) {
	body := &req.CreateDeploymentReqDto{}
	if err := body.Validate(context); err != nil {
		errRes := api_response.BuildErrorResponse(http.StatusBadRequest, "Bad Request", err.Error(), nil)
		context.AbortWithStatusJSON(http.StatusBadRequest, errRes)
		return
	}
	githubRes, code, err := controller.githubHttpClient.ValidateRepositoryByUrl(body.RepositoryUrl)
	if err != nil {
		if code == http.StatusNotFound {
			errRes := api_response.BuildErrorResponse(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "invalid repository", nil)
			context.AbortWithStatusJSON(errRes.Code, errRes)
			return
		}
		errRes := api_response.BuildErrorResponse(code, http.StatusText(code), "", nil)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}
	if githubRes == nil {
		errRes := api_response.BuildErrorResponse(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "invalid repository", nil)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}

	context.JSON(http.StatusOK, githubRes)
}

// DeploymentUpdate
// @Summary  Update Deployment
// @Tags     Deployments
// @Accept   json
// @Produce  json
// @Success  200
// @Router   /deployments/:id [put]
func DeploymentUpdate() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": config.AppConfig.APP_NAME + " is Running",
		})
	}
}

// DeploymentShow
// @Summary  Show Deployment
// @Tags     Deployments
// @Accept   json
// @Produce  json
// @Success  200
// @Router   /deployments/:id [get]
func DeploymentShow() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": config.AppConfig.APP_NAME + " is Running",
		})
	}
}
