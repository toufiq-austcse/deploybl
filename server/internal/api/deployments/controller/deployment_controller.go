package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/toufiq-austcse/deployit/config"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/dto/req"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/mapper"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/service"
	"github.com/toufiq-austcse/deployit/pkg/api_response"
	"github.com/toufiq-austcse/deployit/pkg/http_clients/github"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
)

type DeploymentController struct {
	githubHttpClient  *github.GithubHttpClient
	deploymentService *service.DeploymentService
}

func NewDeploymentController(githubHttpClient *github.GithubHttpClient, deploymentService *service.DeploymentService) *DeploymentController {
	return &DeploymentController{
		githubHttpClient:  githubHttpClient,
		deploymentService: deploymentService,
	}
}

// DeploymentIndex
// @Summary  Deployment Index
// @Tags     Deployments
// @Param        page    query  string  false  "Page"
// @Param        limit   query  string  false  "Limit"
// @Accept   json
// @Produce  json
// @Success  200
// @Router   /deployments [get]
func (controller *DeploymentController) DeploymentIndex() gin.HandlerFunc {
	return func(context *gin.Context) {
		page, _ := strconv.ParseInt(context.DefaultQuery("page", "1"), 10, 64)
		limit, _ := strconv.ParseInt(context.DefaultQuery("limit", "10"), 10, 64)

		deployments, pagination, err := controller.deploymentService.ListDeployment(page, limit, context)
		if err != nil {
			errRes := api_response.BuildErrorResponse(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), err.Error(), nil)
			context.AbortWithStatusJSON(errRes.Code, errRes)
			return
		}

		deploymentListRes := mapper.ToDeploymentListRes(deployments)
		apiRes := api_response.BuildPaginationResponse(http.StatusOK, http.StatusText(http.StatusOK), deploymentListRes, pagination)

		context.JSON(apiRes.Code, apiRes)
	}
}

// DeploymentCreate
// @Summary  Create Deployment
// @Param    request  body  req.CreateDeploymentReqDto  true  "Create Deployment Body"
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
	githubRes, code, err := controller.githubHttpClient.ValidateRepositoryByUrl(&body.RepositoryUrl)
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
	existingDeployment := controller.deploymentService.FindBySubDomainName(&githubRes.Name, context)

	newDeployment := mapper.MapCreateDeploymentReqToSave(body, githubRes, existingDeployment)
	createErr := controller.deploymentService.Create(newDeployment, context)
	if createErr != nil {
		if mongo.IsDuplicateKeyError(createErr) {
			errRes := api_response.BuildErrorResponse(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "this domain name already taken", nil)
			context.AbortWithStatusJSON(errRes.Code, errRes)
			return
		}
		errRes := api_response.BuildErrorResponse(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), createErr.Error(), nil)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}
	createDeploymentRes := api_response.BuildResponse(http.StatusCreated, http.StatusText(http.StatusCreated), mapper.ToDeploymentRes(newDeployment))
	context.JSON(createDeploymentRes.Code, createDeploymentRes)
}

// DeploymentUpdate
// @Summary  Update Deployment
// @Tags     Deployments
// @Param    id  path  string  true  "Deployment ID"
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

// EnvUpdate
// @Summary  Update Deployment Env
// @Tags     Deployments
// @Param    id  path  string  true  "Deployment ID"
// @Accept   json
// @Produce  json
// @Success  200
// @Router   /deployments/:id/env [put]
func (controller *DeploymentController) EnvUpdate(context *gin.Context) {
	var envBody map[string]interface{}
	if err := context.BindJSON(&envBody); err != nil {
		errRes := api_response.BuildErrorResponse(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), err.Error(), nil)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}
	deploymentId := context.Param("id")
	deployment := controller.deploymentService.FindById(deploymentId, context)
	if deployment == nil {
		errRes := api_response.BuildErrorResponse(http.StatusNotFound, http.StatusText(http.StatusNotFound), "", nil)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}
	updatedDeployment, err := controller.deploymentService.UpdateEnv(deploymentId, envBody, context)
	if err != nil {
		errRes := api_response.BuildErrorResponse(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), err.Error(), nil)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}
	deploymentRes := mapper.ToDeploymentDetailsRes(updatedDeployment)

	deploymentDetailsRes := api_response.BuildResponse(http.StatusOK, http.StatusText(http.StatusOK), deploymentRes)
	context.JSON(deploymentDetailsRes.Code, deploymentDetailsRes)
}

// DeploymentShow
// @Summary  Show Deployment
// @Tags     Deployments
// @Param    id  path  string  true  "Deployment ID"
// @Accept   json
// @Produce  json
// @Success  200
// @Router   /deployments/{id} [get]
func (controller *DeploymentController) DeploymentShow(context *gin.Context) {
	deploymentId := context.Param("id")
	deployment := controller.deploymentService.FindById(deploymentId, context)
	if deployment == nil {
		errRes := api_response.BuildErrorResponse(http.StatusNotFound, http.StatusText(http.StatusNotFound), "", nil)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}

	deploymentDetailsRes := api_response.BuildResponse(http.StatusOK, http.StatusText(http.StatusOK), mapper.ToDeploymentDetailsRes(deployment))
	context.JSON(deploymentDetailsRes.Code, deploymentDetailsRes)
}
