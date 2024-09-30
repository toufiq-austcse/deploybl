package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/toufiq-austcse/deployit/enums"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/dto/req"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/mapper"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/service"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/worker"
	"github.com/toufiq-austcse/deployit/pkg/api_response"
	"github.com/toufiq-austcse/deployit/pkg/http_clients/github"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type DeploymentController struct {
	githubHttpClient  *github.GithubHttpClient
	deploymentService *service.DeploymentService
	dockerService     *service.DockerService
	worker            *worker.PullRepoWorker
}

func NewDeploymentController(githubHttpClient *github.GithubHttpClient, deploymentService *service.DeploymentService, repoWorker *worker.PullRepoWorker) *DeploymentController {
	return &DeploymentController{
		githubHttpClient:  githubHttpClient,
		deploymentService: deploymentService,
		worker:            repoWorker,
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
func (controller *DeploymentController) DeploymentIndex(context *gin.Context) {
	page, _ := strconv.ParseInt(context.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(context.DefaultQuery("limit", "10"), 10, 64)
	if page < 1 {
		page = 1
	}

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

	newDeployment := mapper.MapCreateDeploymentReqToSave(body, "github", githubRes, existingDeployment)
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
	controller.worker.PublishPullRepoWork(newDeployment)

	createDeploymentRes := api_response.BuildResponse(http.StatusCreated, http.StatusText(http.StatusCreated), mapper.ToDeploymentRes(newDeployment))
	context.JSON(createDeploymentRes.Code, createDeploymentRes)
}

// DeploymentUpdate
// @Summary  Update Deployment
// @Param    request  body  req.UpdateDeploymentReqDto  true  "Update Deployment Body"
// @Param    id  path  string  true  "Deployment ID"
// @Tags     Deployments
// @Accept   json
// @Produce  json
// @Success  200
// @Router   /deployments/{id} [patch]
func (controller *DeploymentController) DeploymentUpdate(context *gin.Context) {
	deploymentId := context.Param("id")
	body := &req.UpdateDeploymentReqDto{}

	if err := body.Validate(context); err != nil {
		errRes := api_response.BuildErrorResponse(http.StatusBadRequest, "Bad Request", err.Error(), nil)
		context.AbortWithStatusJSON(http.StatusBadRequest, errRes)
		return
	}

	deployment := controller.deploymentService.FindById(deploymentId, context)
	if deployment == nil {
		errRes := api_response.BuildErrorResponse(http.StatusNotFound, http.StatusText(http.StatusNotFound), "", nil)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}
	githubRes, code, err := controller.githubHttpClient.ValidateRepositoryByUrl(&deployment.RepositoryUrl)
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

	updateFields := mapper.MapUpdateDeploymentReqToUpdate(body)
	updatedDeployment, err := controller.deploymentService.UpdateDeployment(deploymentId, updateFields, context)
	if err != nil {
		errRes := api_response.BuildErrorResponse(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), err.Error(), nil)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}
	go func() {
		if updatedDeployment.LatestStatus == enums.QUEUED {
			if deployment.ContainerId != nil {
				fmt.Println("removing old container ", *deployment.ContainerId)
				removeErr := controller.dockerService.RemoveContainer(*deployment.ContainerId)
				if removeErr != nil {
					fmt.Println("error while removing container ", removeErr.Error())
				}
				fmt.Println("container removed ", deployment.ContainerId)
				_, updateErr := controller.deploymentService.UpdateDeployment(deploymentId, map[string]interface{}{
					"container_id": nil,
				}, context)

				if updateErr != nil {
					fmt.Println("error while updating container ", updateErr.Error())
				}

			}
			controller.worker.PublishPullRepoWork(updatedDeployment)
		}
	}()

	deploymentRes := mapper.ToDeploymentDetailsRes(updatedDeployment, *githubRes)
	deploymentDetailsRes := api_response.BuildResponse(http.StatusOK, http.StatusText(http.StatusOK), deploymentRes)
	context.JSON(deploymentDetailsRes.Code, deploymentDetailsRes)

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
	var envBody map[string]string
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
	githubRes, code, err := controller.githubHttpClient.ValidateRepositoryByUrl(&deployment.RepositoryUrl)
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

	updatedDeployment, err := controller.deploymentService.UpdateDeployment(deploymentId, map[string]interface{}{
		"env":           envBody,
		"latest_status": enums.QUEUED,
	}, context)
	if err != nil {
		errRes := api_response.BuildErrorResponse(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), err.Error(), nil)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}
	go func() {
		if updatedDeployment.LatestStatus == enums.QUEUED {
			if deployment.ContainerId != nil {
				fmt.Println("removing old container ", *deployment.ContainerId)
				removeErr := controller.dockerService.RemoveContainer(*deployment.ContainerId)
				if removeErr != nil {
					fmt.Println("error while removing container ", removeErr.Error())
				}
				fmt.Println("container removed ", deployment.ContainerId)
				_, updateErr := controller.deploymentService.UpdateDeployment(deploymentId, map[string]interface{}{
					"container_id": nil,
				}, context)

				if updateErr != nil {
					fmt.Println("error while updating container ", updateErr.Error())
				}

			}
			controller.worker.PublishPullRepoWork(updatedDeployment)
		}
	}()

	deploymentRes := mapper.ToDeploymentDetailsRes(updatedDeployment, *githubRes)

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

	githubRes, code, err := controller.githubHttpClient.ValidateRepositoryByUrl(&deployment.RepositoryUrl)
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

	deploymentDetailsRes := api_response.BuildResponse(http.StatusOK, http.StatusText(http.StatusOK), mapper.ToDeploymentDetailsRes(deployment, *githubRes))
	context.JSON(deploymentDetailsRes.Code, deploymentDetailsRes)
}

func (controller *DeploymentController) DeploymentLatestStatus(context *gin.Context) {
	idsQuery, ok := context.GetQuery("ids")
	if !ok {
		errRes := api_response.BuildErrorResponse(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "ids is required in query param", nil)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}
	idsArray := strings.Split(idsQuery, ",")

	deployments, err := controller.deploymentService.GetLatestStatusByIds(idsArray, context)
	if err != nil {
		errRes := api_response.BuildErrorResponse(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), err.Error(), nil)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}

	deploymentsLatestStatusRes := api_response.BuildResponse(http.StatusOK, http.StatusText(http.StatusOK), mapper.ToDeploymentLatestStatus(deployments))

	context.JSON(deploymentsLatestStatusRes.Code, deploymentsLatestStatusRes)
}
func (controller *DeploymentController) LiveCheckCron(context *gin.Context) {
	runningContainerIds, err := controller.dockerService.ListRunningContainerIds()
	if err != nil {
		errRes := api_response.BuildErrorResponse(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), err.Error(), nil)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}
	updatedCount, err := controller.deploymentService.UpdateDeploymentStatusByContainerIds(runningContainerIds, enums.LIVE, enums.STOPPED, context)
	if err != nil {
		errRes := api_response.BuildErrorResponse(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), err.Error(), nil)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}
	fmt.Println("Time : ", time.Now().Format(time.DateTime), " running containers ", len(runningContainerIds), " stopped ", updatedCount)
	updatedCountRes := api_response.BuildResponse(http.StatusOK, http.StatusText(http.StatusOK), updatedCount)
	context.JSON(updatedCountRes.Code, updatedCountRes)

}
