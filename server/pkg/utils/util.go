package utils

import (
	"context"
	"fmt"
	"os"
	"strings"

	depkloymentModel "github.com/toufiq-austcse/deployit/internal/api/deployments/model"

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

func WriteToFile(text string, event *depkloymentModel.Event) {
	if event == nil {
		return
	}
	createErr := createDirIfNotExists(deployItConfig.AppConfig.EVENT_LOGS_PATH)
	if createErr != nil {
		fmt.Println("error in creating dir ", createErr.Error())
		return
	}
	fileName := deployItConfig.AppConfig.EVENT_LOGS_PATH + "/" + event.Id.Hex() + "_log.txt"
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Println("error in opening file ", err.Error())
	} else {
		defer file.Close()
		if _, writeErr := file.WriteString(text + "\n"); writeErr != nil {
			fmt.Println("error in writing file ", writeErr.Error())
		}
	}
}

func createDirIfNotExists(dirPath string) error {
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0o755)
		if err != nil {
			return err
		}
	}
	return nil
}
