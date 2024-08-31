package app

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/toufiq-austcse/deployit/config"
	"github.com/toufiq-austcse/deployit/di"
	"github.com/toufiq-austcse/deployit/docs"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/controller"
	deploymentRouter "github.com/toufiq-austcse/deployit/internal/api/deployments/router"
	indexRouter "github.com/toufiq-austcse/deployit/internal/api/index/router"
	"github.com/toufiq-austcse/deployit/internal/server"
	"time"
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
	err = container.Invoke(func(deploymentController *controller.DeploymentController) {
		indexRouterGroup := apiServer.GinEngine.Group("")
		indexRouter.Setup(indexRouterGroup)

		deploymentsRouterGroup := apiServer.GinEngine.Group("deployments")
		deploymentRouter.Setup(deploymentsRouterGroup, deploymentController)

	})
	if err != nil {
		return err
	}

	err = apiServer.Run()
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
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Accept", "Origin", "Cache-Control", "X-Requested-With", "Referer", "guest", "publicKey", "Access-Control-Allow-Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	engine.Use(cors.New(corsConfig))
}
