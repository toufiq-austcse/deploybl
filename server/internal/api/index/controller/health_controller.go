package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/toufiq-austcse/deployit/config"
	"github.com/toufiq-austcse/deployit/internal/api/index/event_bus"
	"net/http"
)

type HealthController struct {
	eventBus *event_bus.EventBus
	channel  chan string
}

func NewHealthController(eventBus *event_bus.EventBus) *HealthController {
	ch := make(chan string)
	eventBus.Subscribe("health", ch)
	go HandleEvent(ch)
	return &HealthController{
		eventBus: eventBus,
		channel:  ch,
	}
}

// Index hosts godoc
// @Summary  Health Check
// @Tags     Index
// @Accept   json
// @Produce  json
// @Success  200
// @Router   / [get]
func (healthController *HealthController) Index(context *gin.Context) {
	//err := worker.PublishPullRepoJob()
	//if err != nil {
	//	fmt.Println("error while publishing ", err)
	//}
	healthController.eventBus.Publish("health", "event")

	context.JSON(http.StatusOK, gin.H{
		"message": config.AppConfig.APP_NAME + " is Running",
	})
}

func HandleEvent(eventChan <-chan string) {
	for event := range eventChan {
		fmt.Println("event received: ", event)
	}
}
