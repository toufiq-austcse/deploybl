package di

import (
	_ "github.com/lib/pq" // <------------ here
	"github.com/toufiq-austcse/deployit/internal/api/deployments/controller"
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
		controller.NewDeploymentController,
	}
	for _, provider := range providers {
		if err := c.Provide(provider); err != nil {
			return nil, err
		}
	}
	exec.Command("")
	return c, nil
}
