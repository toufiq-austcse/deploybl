package di

import (
	_ "github.com/lib/pq" // <------------ here
	"github.com/toufiq-austcse/deployit/internal/api/deployments/controller"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/service"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/worker"
	repoController "github.com/toufiq-austcse/deployit/internal/api/repositories/controller"
	userService "github.com/toufiq-austcse/deployit/internal/api/users/service"
	"github.com/toufiq-austcse/deployit/pkg/db/providers/mongodb"
	firebaseClient "github.com/toufiq-austcse/deployit/pkg/firebase"

	"github.com/toufiq-austcse/deployit/pkg/http_clients/github"
	"go.uber.org/dig"
)

func NewDiContainer() (*dig.Container, error) {
	c := dig.New()
	providers := []interface{}{
		mongodb.New,
		github.NewGithubHttpClient,
		firebaseClient.NewFirebaseClient,
		service.NewDeploymentService,
		userService.NewUserService,
		service.NewDockerService,
		controller.NewDeploymentController,
		repoController.NewRepoController,
		worker.NewPullRepoWorker,
		worker.NewBuildRepoWorker,
		worker.NewPreRunRepoWorker,
		worker.NewRunRepoWorker,
		worker.NewStopRepoWorker,
	}
	for _, provider := range providers {
		if err := c.Provide(provider); err != nil {
			return nil, err
		}
	}
	return c, nil
}
