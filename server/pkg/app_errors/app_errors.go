package app_errors

import "errors"

var (
	RepositoryNotFoundError       = errors.New("repository not found")
	DomainNameAlreadyTakenError   = errors.New("domain name already taken")
	DeploymentNotFoundError       = errors.New("deployment not found")
	DeploymentNotRestartableError = errors.New("deployment not restartable")
	DeploymentNotDeployableError  = errors.New("deployment not deployable")
	DeploymentNotStoppableError   = errors.New("deployment not stoppable")
	IdsRequiredInQueryParamError  = errors.New("ids required in query param")
	CannotUpdateError             = errors.New("cannot update")
	DockerImageTagNotFoundError   = errors.New("docker image tag not found")
	ContainerNotFoundError        = errors.New("container not found")
	ContainerPortNotFoundError    = errors.New("container port not found")
)
