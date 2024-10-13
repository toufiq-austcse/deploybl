package mapper

import (
	"time"

	"github.com/sqids/sqids-go"
	"github.com/toufiq-austcse/deployit/config"
	"github.com/toufiq-austcse/deployit/enums"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/dto/req"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/dto/res"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/model"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/worker/payloads"
	userModel "github.com/toufiq-austcse/deployit/internal/api/users/model"
	"github.com/toufiq-austcse/deployit/pkg/http_clients/github/api_res"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MapCreateDeploymentReqToSave(
	dto *req.CreateDeploymentReqDto,
	provider string,
	githubRes *api_res.GithubRepoRes,
	deploymentCount int64,
	user *userModel.User,
) *model.Deployment {
	subDomainName := githubRes.Name
	if deploymentCount > 0 {
		subDomainName += "-" + GenerateShortId(deploymentCount)
	}

	return &model.Deployment{
		Id:                 primitive.NewObjectID(),
		UserId:             user.Id,
		Title:              dto.Title,
		SubDomainName:      subDomainName,
		LatestStatus:       enums.QUEUED,
		LastDeployedAt:     nil,
		RepositoryProvider: provider,
		RepositoryUrl:      dto.RepositoryUrl,
		RepositoryName:     githubRes.Name,
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

func MapUpdateDeploymentReqToUpdate(
	dto *req.UpdateDeploymentReqDto,
	repoName string,
) map[string]interface{} {
	updateFields := map[string]interface{}{}

	if dto.Title != nil {
		updateFields["title"] = *dto.Title
	}
	if dto.BranchName != nil {
		updateFields["branch_name"] = *dto.BranchName
	}
	if dto.RootDir != nil {
		updateFields["root_dir"] = *dto.RootDir
	}
	if dto.DockerFilePath != nil {
		updateFields["docker_file_path"] = *dto.DockerFilePath
	}
	updateFields["repository_name"] = repoName
	if ShouldRedeploy(updateFields) {
		updateFields["latest_status"] = enums.QUEUED
		updateFields["last_deployed_at"] = nil
	}
	return updateFields
}

func ToDeploymentRes(model *model.Deployment) res.DeploymentRes {
	deploymentRes := res.DeploymentRes{
		Id:                 model.Id,
		Title:              model.Title,
		LatestStatus:       model.LatestStatus,
		LastDeployedAt:     model.LastDeployedAt,
		DomainUrl:          GetDomainUrl(model.SubDomainName),
		RepositoryProvider: model.RepositoryProvider,
		BranchName:         model.BranchName,
		CreatedAt:          model.CreatedAt,
		UpdatedAt:          model.UpdatedAt,
	}
	return deploymentRes
}

func ToDeploymentDetailsRes(model *model.Deployment) res.DeploymentDetailsRes {
	deploymentDetail := res.DeploymentDetailsRes{
		Id:                 model.Id,
		Title:              model.Title,
		RepositoryName:     model.RepositoryName,
		DomainUrl:          GetDomainUrl(model.SubDomainName),
		LatestStatus:       model.LatestStatus,
		LastDeployedAt:     model.LastDeployedAt,
		RepositoryProvider: model.RepositoryProvider,
		RepositoryUrl:      model.RepositoryUrl,
		BranchName:         model.BranchName,
		RootDirectory:      model.RootDirectory,
		DockerFilePath:     model.DockerFilePath,
		DockerImageTag:     model.DockerImageTag,
		ContainerId:        model.ContainerId,
		Env:                model.Env,
		CreatedAt:          model.CreatedAt,
		UpdatedAt:          model.UpdatedAt,
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

func ToBuildRepoWorkerPayload(
	payload payloads.PullRepoWorkerPayload,
) payloads.BuildRepoWorkerPayload {
	return payloads.BuildRepoWorkerPayload{
		DeploymentId:   payload.DeploymentId,
		SubDomainName:  payload.SubDomainName,
		DockerFilePath: payload.DockerFilePath,
		BranchName:     payload.BranchName,
		Env:            payload.Env,
	}
}

func ToRunRepoWorkerPayload(payload payloads.BuildRepoWorkerPayload) payloads.RunRepoWorkerPayload {
	return payloads.RunRepoWorkerPayload{
		DeploymentId: payload.DeploymentId,
	}
}

func ToRunRepoWorkerPayloadFromDeployment(
	deployment model.Deployment,
) payloads.RunRepoWorkerPayload {
	return payloads.RunRepoWorkerPayload{
		DeploymentId: deployment.Id.Hex(),
	}
}

func ToStopRepoWorkerPayload(deployment model.Deployment) payloads.StopRepoWorkerPayload {
	return payloads.StopRepoWorkerPayload{DeploymentId: deployment.Id.Hex()}
}

func GetDomainUrl(subDomainName string) string {
	return "https://" + subDomainName + "." + config.AppConfig.BASE_DOMAIN
}

func ToDeploymentLatestStatus(deployments []*model.Deployment) []res.DeploymentLatestStatusRes {
	deploymentsLatestStatusRes := []res.DeploymentLatestStatusRes{}

	for _, deployment := range deployments {
		deploymentsLatestStatusRes = append(
			deploymentsLatestStatusRes,
			res.DeploymentLatestStatusRes{
				Id:             deployment.Id,
				LatestStatus:   deployment.LatestStatus,
				LastDeployedAt: deployment.LastDeployedAt,
				DomainUrl:      GetDomainUrl(deployment.SubDomainName),
			},
		)
	}
	return deploymentsLatestStatusRes
}

func ShouldRedeploy(updatedFieldMap map[string]interface{}) bool {
	for k := range updatedFieldMap {
		if k == "branch_name" || k == "root_dir" || k == "docker_file_path" {
			return true
		}
	}
	return false
}

func GenerateShortId(id int64) string {
	s, _ := sqids.New(sqids.Options{
		MinLength: 5,
		Alphabet:  "abcdefghijklmnopqrstuvwxyz0123456789",
	})
	generatedId, _ := s.Encode([]uint64{uint64(id)})
	return generatedId
}
