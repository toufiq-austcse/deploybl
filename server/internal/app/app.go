package app

import (
	"fmt"
	"net/http"
	"time"

	auth "github.com/toufiq-austcse/deployit/internal/api/auth/controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/toufiq-austcse/deployit/config"
	"github.com/toufiq-austcse/deployit/di"
	"github.com/toufiq-austcse/deployit/docs"
	authRouter "github.com/toufiq-austcse/deployit/internal/api/auth/router"
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

func Run() error {
	if err := config.Init(); err != nil {
		return err
	}
	apiServer := server.NewServer()
	enableCors(apiServer.GinEngine)
	apiServer.GinEngine.Use(func(c *gin.Context) {
		defer func() {
			fmt.Println("recovering")
			if err := recover(); err != nil {
				var errStr string
				switch v := err.(type) {
				case string:
					errStr = v
				case error:
					errStr = v.Error()
				default:
					errStr = fmt.Sprintf("recovered from: %v", v)
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": errStr})
				return
			}
		}()
		c.Next()
	})
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
		eventController *controller.EventController,
		authController *auth.AuthController,
		subscriber *worker.PullRepoWorker,
		firebaseClient *firebase.Client,
		userService *service.UserService,
	) {
		indexRouterGroup := apiServer.GinEngine.Group("")
		indexRouter.Setup(indexRouterGroup)

		authRouterGroup := apiServer.GinEngine.Group("api/v1/auth")
		authRouter.Setup(authRouterGroup, authController)

		deploymentsRouterGroup := apiServer.GinEngine.Group("api/v1/deployments")
		deploymentsRouterGroup.Use(middleware.AuthMiddleware(firebaseClient, userService))
		deploymentRouter.Setup(deploymentsRouterGroup, deploymentController, eventController)

		deploymentCronRouterGroup := apiServer.GinEngine.Group("api/v1/deployments")
		deploymentRouter.SetupDeploymentCronRouter(deploymentCronRouterGroup, deploymentController)

		repositoriesRouterGroup := apiServer.GinEngine.Group("api/v1/repositories")
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
			preRunRepoWorker *worker.PreRunRepoWorker,
		) {
			pullRepoWorker.InitPullRepoSubscriber()
			buildRepoWorker.InitBuildRepoSubscriber()
			preRunRepoWorker.InitPreRunRepoSubscriber()
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
