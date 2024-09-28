package di

import (
	_ "github.com/lib/pq" // <------------ here
	"github.com/toufiq-austcse/deployit/internal/api/deployments/controller"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/service"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/worker"
	repoController "github.com/toufiq-austcse/deployit/internal/api/repositories/controller"
	"github.com/toufiq-austcse/deployit/pkg/db/providers/mongodb"
	"github.com/toufiq-austcse/deployit/pkg/http_clients/github"
	"go.uber.org/dig"
	"os/exec"
)

func NewDiContainer() (*dig.Container, error) {
	c := dig.New()
	providers := []interface {
	}{
		mongodb.New,
		github.NewGithubHttpClient,
		service.NewDeploymentService,
		service.NewDockerService,
		controller.NewDeploymentController,
		repoController.NewRepoController,
		worker.NewPullRepoWorker,
		worker.NewBuildRepoWorker,
		worker.NewRunRepoWorker,
	}
	for _, provider := range providers {
		if err := c.Provide(provider); err != nil {
			return nil, err
		}
	}
	exec.Command("")
	return c, nil
}
