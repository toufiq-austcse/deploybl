package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/toufiq-austcse/deployit/enums/deployment_events_triggered_by"
	"github.com/toufiq-austcse/deployit/enums/reasons"

	deployment_events_enums "github.com/toufiq-austcse/deployit/enums/deployment_events"

	"github.com/gin-gonic/gin"
	"github.com/toufiq-austcse/deployit/enums"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/dto/req"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/mapper"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/service"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/worker"
	"github.com/toufiq-austcse/deployit/pkg/api_response"
	"github.com/toufiq-austcse/deployit/pkg/app_errors"
	"github.com/toufiq-austcse/deployit/pkg/http_clients/github"
	"github.com/toufiq-austcse/deployit/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type DeploymentController struct {
	githubHttpClient  *github.GithubHttpClient
	deploymentService *service.DeploymentService
	dockerService     *service.DockerService
	eventService      *service.EventService
	pullRepoWorker    *worker.PullRepoWorker
	runRepoWorker     *worker.RunRepoWorker
	stopRepoWorker    *worker.StopRepoWorker
	preRunRepoWorker  *worker.PreRunRepoWorker
}

func NewDeploymentController(
	githubHttpClient *github.GithubHttpClient,
	deploymentService *service.DeploymentService,
	eventService *service.EventService,
	pullRepoWorker *worker.PullRepoWorker,
	runRepoWorker *worker.RunRepoWorker,
	stopRepoWorker *worker.StopRepoWorker,
	preRunRepoWorker *worker.PreRunRepoWorker,
) *DeploymentController {
	return &DeploymentController{
		githubHttpClient:  githubHttpClient,
		deploymentService: deploymentService,
		eventService:      eventService,
		pullRepoWorker:    pullRepoWorker,
		runRepoWorker:     runRepoWorker,
		stopRepoWorker:    stopRepoWorker,
		preRunRepoWorker:  preRunRepoWorker,
	}
}

// DeploymentIndex
// @Summary  Deployment Index
// @Tags     Deployments
// @Param    page   query  string  false  "Page"
// @Param    limit  query  string  false  "Limit"
// @Accept   json
// @Produce  json
// @Success  200
// @Router   /api/v1/deployments [get]
// @Success  200  {object}  api_response.Response{data=[]res.DeploymentRes}
func (controller *DeploymentController) DeploymentIndex(context *gin.Context) {
	user := utils.GetUserFromContext(context)

	page, _ := strconv.ParseInt(context.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(context.DefaultQuery("limit", "10"), 10, 64)
	if page < 1 {
		page = 1
	}

	deployments, pagination, err := controller.deploymentService.ListDeployment(
		page,
		limit,
		user.Id,
		context,
	)
	if err != nil {
		errRes := api_response.BuildErrorResponse(
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			err.Error(),
			nil,
		)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}

	deploymentListRes := mapper.ToDeploymentListRes(deployments)
	apiRes := api_response.BuildPaginationResponse(
		http.StatusOK,
		http.StatusText(http.StatusOK),
		deploymentListRes,
		pagination,
	)

	context.JSON(apiRes.Code, apiRes)
}

// DeploymentCreate
// @Summary  Create Deployment
// @Param    request  body  req.CreateDeploymentReqDto  true  "Create Deployment Body"
// @Tags     Deployments
// @Accept   json
// @Produce  json
// @Success  201
// @Router   /api/v1/deployments [post]
// @Success  201  {object}  api_response.Response{data=res.DeploymentRes}
func (controller *DeploymentController) DeploymentCreate(context *gin.Context) {
	user := utils.GetUserFromContext(context)
	body := &req.CreateDeploymentReqDto{}
	if err := body.Validate(context); err != nil {
		errRes := api_response.BuildErrorResponse(
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
			err.Error(),
			nil,
		)
		context.AbortWithStatusJSON(http.StatusBadRequest, errRes)
		return
	}
	githubRes, code, err := controller.githubHttpClient.ValidateRepositoryByUrl(&body.RepositoryUrl)
	if err != nil {
		if code == http.StatusNotFound {
			errRes := api_response.BuildErrorResponse(
				http.StatusBadRequest,
				http.StatusText(http.StatusBadRequest),
				app_errors.RepositoryNotFoundError.Error(),
				nil,
			)
			context.AbortWithStatusJSON(errRes.Code, errRes)
			return
		}
		errRes := api_response.BuildErrorResponse(
			code,
			http.StatusText(code),
			http.StatusText(code),
			nil,
		)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}
	if githubRes == nil {
		errRes := api_response.BuildErrorResponse(
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
			app_errors.RepositoryNotFoundError.Error(),
			nil,
		)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}

	deploymentCount, _ := controller.deploymentService.CountDeploymentByRepositoryName(
		githubRes.Name,
		context,
	)

	newDeployment := mapper.MapCreateDeploymentReqToSave(
		body,
		"github",
		githubRes,
		deploymentCount,
		user,
	)
	event, createErr := controller.deploymentService.Create(newDeployment, context)
	if createErr != nil {
		if mongo.IsDuplicateKeyError(createErr) {
			errRes := api_response.BuildErrorResponse(
				http.StatusBadRequest,
				http.StatusText(http.StatusBadRequest),
				app_errors.DomainNameAlreadyTakenError.Error(),
				nil,
			)
			context.AbortWithStatusJSON(errRes.Code, errRes)
			return
		}
		errRes := api_response.BuildErrorResponse(
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			createErr.Error(),
			nil,
		)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}
	controller.pullRepoWorker.PublishPullRepoWork(newDeployment, event)

	createDeploymentRes := api_response.BuildResponse(
		http.StatusCreated,
		http.StatusText(http.StatusCreated),
		mapper.ToDeploymentRes(newDeployment),
	)
	context.JSON(createDeploymentRes.Code, createDeploymentRes)
}

// DeploymentUpdate
// @Summary  Update Deployment
// @Param    request  body  req.UpdateDeploymentReqDto  true  "Update Deployment Body"
// @Param    id       path  string                      true  "Deployment ID"
// @Tags     Deployments
// @Accept   json
// @Produce  json
// @Success  200
// @Router   /api/v1/deployments/{id} [patch]
// @Success  200  {object}  api_response.Response{data=res.DeploymentDetailsRes}
func (controller *DeploymentController) DeploymentUpdate(context *gin.Context) {
	deploymentId := context.Param("id")
	user := utils.GetUserFromContext(context)
	body := &req.UpdateDeploymentReqDto{}

	if err := body.Validate(context); err != nil {
		errRes := api_response.BuildErrorResponse(
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
			err.Error(),
			nil,
		)
		context.AbortWithStatusJSON(http.StatusBadRequest, errRes)
		return
	}

	deployment := controller.deploymentService.FindUserDeploymentById(
		deploymentId,
		user.Id,
		context,
	)
	if deployment == nil {
		errRes := api_response.BuildErrorResponse(
			http.StatusNotFound,
			http.StatusText(http.StatusNotFound),
			app_errors.DeploymentNotFoundError.Error(),
			nil,
		)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}
	githubRes, code, err := controller.githubHttpClient.ValidateRepositoryByUrl(
		&deployment.RepositoryUrl,
	)
	if err != nil {
		if code == http.StatusNotFound {
			errRes := api_response.BuildErrorResponse(
				http.StatusBadRequest,
				http.StatusText(http.StatusBadRequest),
				app_errors.RepositoryNotFoundError.Error(),
				nil,
			)
			context.AbortWithStatusJSON(errRes.Code, errRes)
			return
		}
		errRes := api_response.BuildErrorResponse(
			code,
			http.StatusText(code),
			http.StatusText(code),
			nil,
		)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}
	if githubRes == nil {
		errRes := api_response.BuildErrorResponse(
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
			app_errors.RepositoryNotFoundError.Error(),
			nil,
		)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}

	updateFields := mapper.MapUpdateDeploymentReqToUpdate(body, githubRes.FullName)
	updatedDeployment, err := controller.deploymentService.UpdateDeployment(
		deploymentId,
		updateFields,
		nil,
		context,
	)
	if err != nil {
		errRes := api_response.BuildErrorResponse(
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			err.Error(),
			nil,
		)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}

	newEvent := mapper.MapEventModelToSave(
		deployment.Id,
		deployment_events_enums.REBUILD_DEPLOYMENT,
		deployment_events_triggered_by.USER,
		user.Id.Hex(),
		reasons.GetReasonPtr(reasons.SETTINGS_UPDATE),
	)
	eventErr := controller.eventService.Create(newEvent, context)
	if eventErr != nil {
		fmt.Println("error in creating event ", eventErr.Error())
	}
	go func() {
		if updatedDeployment.LatestStatus == enums.QUEUED {
			controller.pullRepoWorker.PublishPullRepoWork(updatedDeployment, newEvent)
		}
	}()

	deploymentRes := mapper.ToDeploymentDetailsRes(updatedDeployment)
	deploymentDetailsRes := api_response.BuildResponse(
		http.StatusOK,
		http.StatusText(http.StatusOK),
		deploymentRes,
	)
	context.JSON(deploymentDetailsRes.Code, deploymentDetailsRes)
}

// EnvUpdate
// @Summary  Update Deployment Env
// @Tags     Deployments
// @Param    id  path  string  true  "Deployment ID"
// @Accept   json
// @Produce  json
// @Success  200
// @Router   /api/v1/deployments/{id}/env [put]
// @Success  200  {object}  api_response.Response{data=res.DeploymentDetailsRes}
func (controller *DeploymentController) EnvUpdate(context *gin.Context) {
	var envBody map[string]string
	if err := context.BindJSON(&envBody); err != nil {
		errRes := api_response.BuildErrorResponse(
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
			err.Error(),
			nil,
		)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}
	deploymentId := context.Param("id")
	user := utils.GetUserFromContext(context)
	deployment := controller.deploymentService.FindUserDeploymentById(
		deploymentId,
		user.Id,
		context,
	)
	if deployment == nil {
		errRes := api_response.BuildErrorResponse(
			http.StatusNotFound,
			http.StatusText(http.StatusNotFound),
			app_errors.DeploymentNotFoundError.Error(),
			nil,
		)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}
	githubRes, code, err := controller.githubHttpClient.ValidateRepositoryByUrl(
		&deployment.RepositoryUrl,
	)
	if err != nil {
		if code == http.StatusNotFound {
			errRes := api_response.BuildErrorResponse(
				http.StatusBadRequest,
				http.StatusText(http.StatusBadRequest),
				app_errors.RepositoryNotFoundError.Error(),
				nil,
			)
			context.AbortWithStatusJSON(errRes.Code, errRes)
			return
		}
		errRes := api_response.BuildErrorResponse(
			code,
			http.StatusText(code),
			http.StatusText(code),
			nil,
		)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}
	if githubRes == nil {
		errRes := api_response.BuildErrorResponse(
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
			app_errors.RepositoryNotFoundError.Error(),
			nil,
		)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}

	updatedDeployment, err := controller.deploymentService.UpdateDeployment(
		deploymentId,
		map[string]interface{}{
			"env":           envBody,
			"latest_status": enums.QUEUED,
		},
		nil,
		context,
	)
	if err != nil {
		errRes := api_response.BuildErrorResponse(
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			err.Error(),
			nil,
		)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}
	newEvent := mapper.MapEventModelToSave(
		deployment.Id,
		deployment_events_enums.RESTART_DEPLOYMENT,
		deployment_events_triggered_by.USER, user.Id.Hex(),
		reasons.GetReasonPtr(reasons.ENV_UPDATED))
	if eventErr := controller.eventService.Create(newEvent, context); eventErr != nil {
		fmt.Println("error in creating event ", eventErr.Error())
	}
	go func() {
		runRepoJobPayload := mapper.ToRunRepoWorkerPayloadFromDeployment(*deployment, *newEvent)
		publishJobErr := controller.runRepoWorker.PublishRunRepoJob(runRepoJobPayload)
		if publishJobErr != nil {
			fmt.Println("error while publishing job ", publishJobErr.Error())
		}
	}()

	deploymentRes := mapper.ToDeploymentDetailsRes(updatedDeployment)

	deploymentDetailsRes := api_response.BuildResponse(
		http.StatusOK,
		http.StatusText(http.StatusOK),
		deploymentRes,
	)
	context.JSON(deploymentDetailsRes.Code, deploymentDetailsRes)
}

// DeploymentShow
// @Summary  Show Deployment
// @Tags     Deployments
// @Param    id  path  string  true  "Deployment ID"
// @Accept   json
// @Produce  json
// @Success  200
// @Router   /api/v1/deployments/{id} [get]
// @Success  200  {object}  api_response.Response{data=res.DeploymentDetailsRes}
func (controller *DeploymentController) DeploymentShow(context *gin.Context) {
	deploymentId := context.Param("id")
	user := utils.GetUserFromContext(context)
	deployment := controller.deploymentService.FindUserDeploymentById(
		deploymentId,
		user.Id,
		context,
	)
	if deployment == nil {
		errRes := api_response.BuildErrorResponse(
			http.StatusNotFound,
			http.StatusText(http.StatusNotFound),
			app_errors.DeploymentNotFoundError.Error(),
			nil,
		)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}

	deploymentDetailsRes := api_response.BuildResponse(
		http.StatusOK,
		http.StatusText(http.StatusOK),
		mapper.ToDeploymentDetailsRes(deployment),
	)
	context.JSON(deploymentDetailsRes.Code, deploymentDetailsRes)
}

// DeploymentRestart
// @Summary  Restart Deployment
// @Tags     Deployments
// @Param    id  path  string  true  "Deployment ID"
// @Accept   json
// @Produce  json
// @Success  200
// @Router   /api/v1/deployments/{id}/restart [post]
// @Success  200  {object}  api_response.Response{data=res.DeploymentDetailsRes}
func (controller *DeploymentController) DeploymentRestart(context *gin.Context) {
	deploymentId := context.Param("id")
	user := utils.GetUserFromContext(context)
	deployment := controller.deploymentService.FindUserDeploymentById(
		deploymentId,
		user.Id,
		context,
	)
	if deployment == nil {
		errRes := api_response.BuildErrorResponse(
			http.StatusNotFound,
			http.StatusText(http.StatusNotFound),
			app_errors.DeploymentNotFoundError.Error(),
			nil,
		)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}
	if !controller.deploymentService.IsRestartable(deployment) {
		errRes := api_response.BuildErrorResponse(
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
			app_errors.DeploymentNotRestartableError.Error(),
			nil,
		)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}
	updatedDeployment, err := controller.deploymentService.UpdateLatestStatus(
		deploymentId,
		enums.QUEUED,
		nil,
		context,
	)
	if err != nil {
		errRes := api_response.BuildErrorResponse(
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			err.Error(),
			nil,
		)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}
	newEvent := mapper.MapEventModelToSave(
		deployment.Id,
		deployment_events_enums.RESTART_DEPLOYMENT,
		deployment_events_triggered_by.USER,
		user.Id.Hex(),
		reasons.GetReasonPtr(reasons.TRIGGERED_VIA_DASHBOARD),
	)

	eventErr := controller.eventService.Create(newEvent, context)
	if eventErr != nil {
		fmt.Println("error in creating event ", eventErr.Error())
	}

	go func() {
		runRepoJobPayload := mapper.ToPreRunRepoFromDeployment(*deployment, *newEvent)
		publishJobErr := controller.preRunRepoWorker.PublishPreRunRepoJob(runRepoJobPayload)

		if publishJobErr != nil {
			fmt.Println("error while publishing job ", publishJobErr.Error())
		}
	}()

	deploymentRes := mapper.ToDeploymentDetailsRes(updatedDeployment)

	deploymentDetailsRes := api_response.BuildResponse(
		http.StatusOK,
		http.StatusText(http.StatusOK),
		deploymentRes,
	)
	context.JSON(deploymentDetailsRes.Code, deploymentDetailsRes)
}

// DeploymentRebuildAndReDeploy
// @Summary  Rebuild and Deploy Deployment
// @Tags     Deployments
// @Param    id  path  string  true  "Deployment ID"
// @Accept   json
// @Produce  json
// @Success  200
// @Router   /api/v1/deployments/{id}/rebuild-and-redeploy [post]
// @Success  200  {object}  api_response.Response{data=res.DeploymentDetailsRes}
func (controller *DeploymentController) DeploymentRebuildAndReDeploy(context *gin.Context) {
	deploymentId := context.Param("id")
	user := utils.GetUserFromContext(context)
	deployment := controller.deploymentService.FindUserDeploymentById(
		deploymentId,
		user.Id,
		context,
	)
	if deployment == nil {
		errRes := api_response.BuildErrorResponse(
			http.StatusNotFound,
			http.StatusText(http.StatusNotFound),
			app_errors.RepositoryNotFoundError.Error(),
			nil,
		)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}
	if !controller.deploymentService.IsRebuildAble(deployment) {
		errRes := api_response.BuildErrorResponse(
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
			app_errors.DeploymentNotDeployableError.Error(),
			nil,
		)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}
	updatedDeployment, err := controller.deploymentService.UpdateLatestStatus(
		deploymentId,
		enums.QUEUED,
		nil,
		context,
	)
	if err != nil {
		errRes := api_response.BuildErrorResponse(
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			err.Error(),
			nil,
		)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}
	newEvent := mapper.MapEventModelToSave(
		deployment.Id,
		deployment_events_enums.REBUILD_DEPLOYMENT,
		deployment_events_triggered_by.USER,
		user.Id.Hex(),
		reasons.GetReasonPtr(reasons.TRIGGERED_VIA_DASHBOARD),
	)
	eventErr := controller.eventService.Create(newEvent, context)
	if eventErr != nil {
		fmt.Println("error in creating event ", eventErr.Error())
	}
	go func() {
		controller.pullRepoWorker.PublishPullRepoWork(updatedDeployment, newEvent)
	}()

	deploymentRes := mapper.ToDeploymentDetailsRes(updatedDeployment)

	deploymentDetailsRes := api_response.BuildResponse(
		http.StatusOK,
		http.StatusText(http.StatusOK),
		deploymentRes,
	)
	context.JSON(deploymentDetailsRes.Code, deploymentDetailsRes)
}

// DeploymentStop
// @Summary  Stop Deployment
// @Tags     Deployments
// @Param    id  path  string  true  "Deployment ID"
// @Accept   json
// @Produce  json
// @Success  200
// @Router   /api/v1/deployments/{id}/stop [post]
// @Success  200  {object}  api_response.Response{data=res.DeploymentDetailsRes}
func (controller *DeploymentController) DeploymentStop(context *gin.Context) {
	deploymentId := context.Param("id")
	user := utils.GetUserFromContext(context)
	deployment := controller.deploymentService.FindUserDeploymentById(
		deploymentId,
		user.Id,
		context,
	)
	if deployment == nil {
		errRes := api_response.BuildErrorResponse(
			http.StatusNotFound,
			http.StatusText(http.StatusNotFound),
			app_errors.DeploymentNotFoundError.Error(),
			nil,
		)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}
	if !controller.deploymentService.IsStopAble(deployment) {
		errRes := api_response.BuildErrorResponse(
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
			app_errors.DeploymentNotStoppableError.Error(),
			nil,
		)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}

	updatedDeployment, err := controller.deploymentService.UpdateLatestStatus(
		deploymentId,
		enums.QUEUED,
		nil,
		context,
	)
	if err != nil {
		errRes := api_response.BuildErrorResponse(
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			err.Error(),
			nil,
		)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}

	newEvent := mapper.MapEventModelToSave(
		deployment.Id,
		deployment_events_enums.STOP_DEPLOYMENT,
		deployment_events_triggered_by.USER,
		user.Id.Hex(),
		reasons.GetReasonPtr(reasons.TRIGGERED_VIA_DASHBOARD),
	)
	eventErr := controller.eventService.Create(newEvent, context)
	if eventErr != nil {
		fmt.Println("error in creating event ", eventErr.Error())
	}

	stopRepoWorkerPayload := mapper.ToStopRepoWorkerPayload(*deployment, newEvent)

	if publishErr := controller.stopRepoWorker.PublishStopRepoJob(stopRepoWorkerPayload); publishErr != nil {
		errRes := api_response.BuildErrorResponse(
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			publishErr.Error(),
			nil,
		)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}

	deploymentRes := mapper.ToDeploymentDetailsRes(updatedDeployment)

	deploymentDetailsRes := api_response.BuildResponse(
		http.StatusOK,
		http.StatusText(http.StatusOK),
		deploymentRes,
	)
	context.JSON(deploymentDetailsRes.Code, deploymentDetailsRes)
}

// DeploymentLatestStatus
// @Summary  Deployments Latest Status
// @Tags     Deployments
// @Param    ids  query  string  true  "Deployment ID"
// @Accept   json
// @Produce  json
// @Success  200
// @Router   /api/v1/deployments/latest-status [get]
// @Success  200  {object}  api_response.Response{data=[]res.DeploymentLatestStatusRes}
func (controller *DeploymentController) DeploymentLatestStatus(context *gin.Context) {
	idsQuery, ok := context.GetQuery("ids")
	user := utils.GetUserFromContext(context)
	if !ok {
		errRes := api_response.BuildErrorResponse(
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
			app_errors.IdsRequiredInQueryParamError.Error(),
			nil,
		)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}
	idsArray := strings.Split(idsQuery, ",")
	if len(idsArray) == 1 && idsArray[0] == "" {
		errRes := api_response.BuildErrorResponse(
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
			app_errors.IdsRequiredInQueryParamError.Error(),
			nil,
		)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}

	deployments, err := controller.deploymentService.GetLatestStatusByIds(
		idsArray,
		user.Id,
		context,
	)
	if err != nil {
		errRes := api_response.BuildErrorResponse(
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			err.Error(),
			nil,
		)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}

	deploymentsLatestStatusRes := api_response.BuildResponse(
		http.StatusOK,
		http.StatusText(http.StatusOK),
		mapper.ToDeploymentLatestStatus(deployments),
	)

	context.JSON(deploymentsLatestStatusRes.Code, deploymentsLatestStatusRes)
}

// LiveCheckCron
// @Summary  Check Stopped Deployments
// @Tags     Deployments
// @Accept   json
// @Produce  json
// @Success  200
// @Router   /api/v1/deployments/check-stopped-cron [post]
// @Success  200  {object}  api_response.Response{data=int}
func (controller *DeploymentController) LiveCheckCron(context *gin.Context) {
	runningContainerIds, err := controller.dockerService.ListRunningContainerIds()
	if err != nil {
		errRes := api_response.BuildErrorResponse(
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			err.Error(),
			nil,
		)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}
	updatedCount, err := controller.deploymentService.UpdateDeploymentStatusByContainerIds(
		runningContainerIds,
		enums.LIVE,
		enums.STOPPED,
		context,
	)
	if err != nil {
		errRes := api_response.BuildErrorResponse(
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			err.Error(),
			nil,
		)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}
	fmt.Println(
		"Time : ",
		time.Now().Format(time.DateTime),
		" running containers ",
		len(runningContainerIds),
		" stopped ",
		updatedCount,
	)
	updatedCountRes := api_response.BuildResponse(
		http.StatusOK,
		http.StatusText(http.StatusOK),
		updatedCount,
	)
	context.JSON(updatedCountRes.Code, updatedCountRes)
}

// DeployingCheckCron
// @Summary  Check Deploying state Deployments
// @Tags     Deployments
// @Accept   json
// @Produce  json
// @Success  200
// @Router   /api/v1/deployments/check-deploying-cron [post]
// @Success  200  {object}  api_response.Response{data=int}
func (controller *DeploymentController) DeployingCheckCron(context *gin.Context) {
	count := 0

	events, err := controller.eventService.GetLatestProcessingEventsByDeployments(context)
	if err != nil {
		errRes := api_response.BuildErrorResponse(
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			err.Error(),
			nil,
		)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}

	deploymentIds := controller.eventService.GetDeploymentIdFromEvents(events)
	deployments, err := controller.deploymentService.GetDeploymentsByIds(deploymentIds, enums.DEPLOYING, context)
	if err != nil {
		errRes := api_response.BuildErrorResponse(
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			err.Error(),
			nil,
		)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}

	for _, deployment := range deployments {
		event := controller.eventService.FindEventByDeploymentId(events, deployment.Id)
		if event == nil {
			continue
		}
		runRepoPayload := mapper.ToRunRepoWorkerPayloadFromDeployment(deployment, *event)
		if publishErr := controller.runRepoWorker.PublishRunRepoJob(runRepoPayload); publishErr != nil {
			fmt.Println("error while publishing job ", publishErr.Error())
		} else {
			count++
		}
	}
	cronRes := api_response.BuildResponse(
		http.StatusOK,
		http.StatusText(http.StatusOK),
		count,
	)
	context.JSON(http.StatusOK, cronRes)
}
