package mapper

import (
	"github.com/toufiq-austcse/deployit/config"
	"github.com/toufiq-austcse/deployit/enums"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/dto/req"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/dto/res"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/model"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/worker/payloads"
	"github.com/toufiq-austcse/deployit/pkg/http_clients/github/api_res"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"time"
)

func MapCreateDeploymentReqToSave(dto *req.CreateDeploymentReqDto, provider string, githubRes *api_res.GithubRepoRes, existingDeployment *model.Deployment) *model.Deployment {
	subDomainName := githubRes.Name
	if existingDeployment != nil {
		subDomainName += "-" + strconv.FormatInt(time.Now().UnixMilli(), 10)
	}

	return &model.Deployment{
		Id:                 primitive.NewObjectID(),
		Title:              dto.Title,
		SubDomainName:      subDomainName,
		LatestStatus:       enums.QUEUED,
		LastDeployedAt:     nil,
		RepositoryProvider: provider,
		RepositoryUrl:      dto.RepositoryUrl,
		GitUrl:             githubRes.CloneURL,
		BranchName:         dto.BranchName,
		RootDirectory:      dto.RootDir,
		DockerFilePath:     *dto.DockerFilePath,
		DockerImageTag:     nil,
		ContainerId:        nil,
		Env:                dto.Env,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
}
func ToDeploymentRes(model *model.Deployment) res.DeploymentRes {
	return res.DeploymentRes{
		Id:                 model.Id,
		Title:              model.Title,
		LatestStatus:       model.LatestStatus,
		LastDeployedAt:     model.LastDeployedAt,
		RepositoryProvider: model.RepositoryProvider,
		BranchName:         model.BranchName,
		CreatedAt:          model.CreatedAt,
		UpdatedAt:          model.UpdatedAt,
	}
}
func ToDeploymentDetailsRes(model *model.Deployment, githubRes api_res.GithubRepoRes) res.DeploymentDetailsRes {
	deploymentDetail := res.DeploymentDetailsRes{
		Id:                 model.Id,
		Title:              model.Title,
		RepositoryName:     githubRes.FullName,
		DomainUrl:          nil,
		LatestStatus:       model.LatestStatus,
		LastDeployedAt:     model.LastDeployedAt,
		RepositoryProvider: model.RepositoryProvider,
		RepositoryUrl:      model.RepositoryUrl,
		BranchName:         model.BranchName,
		DockerFilePath:     model.DockerFilePath,
		DockerImageTag:     model.DockerImageTag,
		ContainerId:        model.ContainerId,
		Env:                model.Env,
		CreatedAt:          model.CreatedAt,
		UpdatedAt:          model.UpdatedAt,
	}
	if model.LatestStatus == enums.LIVE {
		domainUrl := GetDomainUrl(model.SubDomainName)
		deploymentDetail.DomainUrl = &domainUrl
	}
	return deploymentDetail
}

func ToDeploymentListRes(deploymentModels []*model.Deployment) []res.DeploymentRes {
	deployments := []res.DeploymentRes{}

	for _, deployment := range deploymentModels {
		deployments = append(deployments, ToDeploymentRes(deployment))
	}

	return deployments
}
func ToPullRepoWorkerPayload(deployment *model.Deployment) payloads.PullRepoWorkerPayload {
	return payloads.PullRepoWorkerPayload{
		DeploymentId:   deployment.Id.Hex(),
		BranchName:     deployment.BranchName,
		SubDomainName:  deployment.SubDomainName,
		GitUrl:         deployment.GitUrl,
		RootDir:        deployment.RootDirectory,
		DockerFilePath: deployment.DockerFilePath,
		Env:            deployment.Env,
	}
}

func ToBuildRepoWorkerPayload(payload payloads.PullRepoWorkerPayload) payloads.BuildRepoWorkerPayload {
	return payloads.BuildRepoWorkerPayload{
		DeploymentId:   payload.DeploymentId,
		SubDomainName:  payload.SubDomainName,
		DockerFilePath: payload.DockerFilePath,
		RootDir:        payload.RootDir,
		Env:            payload.Env,
	}
}
func ToRunRepoWorkerPayload(payload payloads.BuildRepoWorkerPayload, dockerImageTag string) payloads.RunRepoWorkerPayload {
	return payloads.RunRepoWorkerPayload{
		DeploymentId:   payload.DeploymentId,
		DockerImageTag: dockerImageTag,
		Env:            payload.Env,
	}
}
func GetDomainUrl(subDomainName string) string {
	return "https://" + subDomainName + "." + config.AppConfig.BASE_DOMAIN
}
func ToDeploymentLatestStatus(deployments []*model.Deployment) []res.DeploymentLatestStatusRes {
	deploymentsLatestStatusRes := []res.DeploymentLatestStatusRes{}

	for _, deployment := range deployments {
		var deploymentDomainUrl *string = nil
		if deployment.LatestStatus == enums.LIVE {
			domainUrl := GetDomainUrl(deployment.SubDomainName)
			deploymentDomainUrl = &domainUrl
		}
		deploymentsLatestStatusRes = append(deploymentsLatestStatusRes, res.DeploymentLatestStatusRes{
			Id:             deployment.Id,
			LatestStatus:   deployment.LatestStatus,
			LastDeployedAt: deployment.LastDeployedAt,
			DomainUrl:      deploymentDomainUrl,
		})
	}
	return deploymentsLatestStatusRes

}
