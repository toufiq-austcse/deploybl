package utils

import (
	"context"
	deployItConfig "github.com/toufiq-austcse/deployit/config"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/worker/payloads"
	"github.com/toufiq-austcse/deployit/internal/api/users/model"
	"strings"
)

func GetDockerBuildContextPath(payload payloads.BuildRepoWorkerPayload) string {
	if payload.RootDir == nil {
		return deployItConfig.AppConfig.REPOSITORIES_PATH + "/" + payload.DeploymentId
	}
	return deployItConfig.AppConfig.REPOSITORIES_PATH + "/" + payload.DeploymentId + "/" + *payload.RootDir

}

func GetLocalRepoPath(DeploymentId string, branchName string) string {
	return deployItConfig.AppConfig.REPOSITORIES_PATH + "/" + DeploymentId + "-" + branchName

}

func ParseRepositoryUrl(repoUrl string) string {
	return strings.Split(repoUrl, ".git")[0]
}

func GetUserFromContext(ctx context.Context) *model.User {
	return ctx.Value("user").(*model.User)
}
