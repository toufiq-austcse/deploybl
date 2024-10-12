package app

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/toufiq-austcse/deployit/config"
	"github.com/toufiq-austcse/deployit/di"
	"github.com/toufiq-austcse/deployit/docs"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/controller"
	deploymentRouter "github.com/toufiq-austcse/deployit/internal/api/deployments/router"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/worker"
	indexRouter "github.com/toufiq-austcse/deployit/internal/api/index/router"
	repoController "github.com/toufiq-austcse/deployit/internal/api/repositories/controller"
	repoRouter "github.com/toufiq-austcse/deployit/internal/api/repositories/router"
	"github.com/toufiq-austcse/deployit/internal/api/users/service"
	"github.com/toufiq-austcse/deployit/internal/middleware"
	"github.com/toufiq-austcse/deployit/internal/server"
	"github.com/toufiq-austcse/deployit/pkg/firebase"
	"go.uber.org/dig"
)

func Run(configPath string) error {
	config.Init(configPath)
	apiServer := server.NewServer()
	enableCors(apiServer.GinEngine)
	setupSwagger(apiServer)
	container, err := di.NewDiContainer()
	if err != nil {
		return err
	}
	err = SetupRouters(apiServer, container)
	if err != nil {
		return err
	}
	err = SetupSubscribers(container)
	if err != nil {
		return err
	}

	err = apiServer.Run()
	if err != nil {
		return err
	}

	return nil
}

func SetupRouters(apiServer *server.Server, container *dig.Container) error {
	err := container.Invoke(func(deploymentController *controller.DeploymentController,
		repoController *repoController.RepoController,
		subscriber *worker.PullRepoWorker,
		firebaseClient *firebase.Client,
		userService *service.UserService,
	) {
		indexRouterGroup := apiServer.GinEngine.Group("")
		indexRouter.Setup(indexRouterGroup)

		deploymentsRouterGroup := apiServer.GinEngine.Group("deployments")
		deploymentsRouterGroup.Use(middleware.AuthMiddleware(firebaseClient, userService))
		deploymentRouter.Setup(deploymentsRouterGroup, deploymentController)

		deploymentCronRouterGroup := apiServer.GinEngine.Group("deployments")
		deploymentRouter.SetupDeploymentCronRouter(deploymentCronRouterGroup, deploymentController)

		repositoriesRouterGroup := apiServer.GinEngine.Group("repositories")
		repositoriesRouterGroup.Use(middleware.AuthMiddleware(firebaseClient, userService))
		repoRouter.Setup(repositoriesRouterGroup, repoController)
	})
	if err != nil {
		return err
	}
	return nil
}

func SetupSubscribers(container *dig.Container) error {
	err := container.Invoke(
		func(pullRepoWorker *worker.PullRepoWorker, buildRepoWorker *worker.BuildRepoWorker,
			runRepoWorker *worker.RunRepoWorker, stopRepoWorker *worker.StopRepoWorker,
		) {
			pullRepoWorker.InitPullRepoSubscriber()
			buildRepoWorker.InitBuildRepoSubscriber()
			runRepoWorker.InitRunRepoSubscriber()
			stopRepoWorker.InitStopRepoSubscriber()
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func setupSwagger(apiServer *server.Server) {
	docs.SwaggerInfo.Title = config.AppConfig.APP_NAME + " API DOC"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Host = config.AppConfig.APP_URL
	apiServer.GinEngine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func enableCors(engine *gin.Engine) {
	corsConfig := cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders: []string{
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
			"Accept",
			"Origin",
			"Cache-Control",
			"X-Requested-With",
			"Referer",
			"guest",
			"publicKey",
			"Access-Control-Allow-Origin",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	engine.Use(cors.New(corsConfig))
}
