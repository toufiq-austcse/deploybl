package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/worker"
	"github.com/toufiq-austcse/deployit/internal/api/repositories/mapper"
	"github.com/toufiq-austcse/deployit/pkg/api_response"
	"github.com/toufiq-austcse/deployit/pkg/http_clients/github"
	"net/http"
)

type RepoController struct {
	githubHttpClient *github.GithubHttpClient
	subscriber       *worker.PullRepoWorker
}

func NewRepoController(githubHttpClient *github.GithubHttpClient, subscriber *worker.PullRepoWorker) *RepoController {
	return &RepoController{
		githubHttpClient: githubHttpClient,
		subscriber:       subscriber,
	}
}

// GetRepoDetails
// @Summary  Get Repo Details
// @Tags     Repositories
// @Param    repo_url  query  string  true  "Repo Url"
// @Accept   json
// @Produce  json
// @Success  200
// @Router   /repositories [get]
func (controller RepoController) GetRepoDetails(context *gin.Context) {
	repoUrl, isExist := context.GetQuery("repo_url")
	if !isExist {
		errRes := api_response.BuildErrorResponse(http.StatusBadRequest, "repo_url required in query", "", nil)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}
	githubRes, code, err := controller.githubHttpClient.ValidateRepositoryByUrl(&repoUrl)
	if err != nil {
		errRes := api_response.BuildErrorResponse(code, http.StatusText(code), err.Error(), nil)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}

	detailsRes := mapper.ToRepoDetailsRes(githubRes)

	repoDetailsApiRes := api_response.BuildResponse(http.StatusOK, http.StatusText(http.StatusOK), detailsRes)
	context.JSON(repoDetailsApiRes.Code, repoDetailsApiRes)
}

// PublishMsh
// @Summary  Get Repo Details
// @Tags     Repositories
// @Param    repo_url  query  string  true  "Repo Url"
// @Accept   json
// @Produce  json
// @Success  200
// @Router   /publish [get]
func (controller RepoController) PublishMsh(context *gin.Context) {
	err := controller.subscriber.PublishPullRepoJob()
	if err != nil {
		context.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	context.JSON(http.StatusOK, nil)
}
