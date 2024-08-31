package di

import (
	_ "github.com/lib/pq" // <------------ here
	"github.com/toufiq-austcse/go-api-boilerplate/pkg/db/providers/mongodb"
	"go.uber.org/dig"
	"os/exec"
)

func NewDiContainer() (*dig.Container, error) {
	c := dig.New()
	providers := []interface {
	}{
		mongodb.New,
	}
	for _, provider := range providers {
		if err := c.Provide(provider); err != nil {
			return nil, err
		}
	}
	exec.Command("")
	return c, nil
}
