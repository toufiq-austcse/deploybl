package utils

import deployItConfig "github.com/toufiq-austcse/deployit/config"

func GetLocalRepoPath(deploymentID string) string {
	return deployItConfig.AppConfig.REPOSITORIES_PATH + "/" + deploymentID
}
