package utils

import (
	"context"
	"os"
	"strings"

	deployItConfig "github.com/toufiq-austcse/deployit/config"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/worker/payloads"
	"github.com/toufiq-austcse/deployit/internal/api/users/model"
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

func CreateDirIfNotExists(dirPath string) error {
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0o755)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetEventLogFilePath(eventId string) string {
	return deployItConfig.AppConfig.EVENT_LOGS_PATH + "/" + eventId + "_log.txt"
}
