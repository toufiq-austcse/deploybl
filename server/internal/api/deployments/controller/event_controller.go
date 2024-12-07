package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/mapper"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/service"
	"github.com/toufiq-austcse/deployit/pkg/api_response"
)

type EventController struct {
	eventService *service.EventService
}

func NewEventController(eventService *service.EventService) *EventController {
	return &EventController{
		eventService: eventService,
	}
}

// EventIndex
// @Summary  Deployment Events
// @Tags     Deployments
// @Param    page   query  string  false  "Page"
// @Param    limit  query  string  false  "Limit"
// @Param    id     path   string  true   "Deployment ID"
// @Accept   json
// @Produce  json
// @Success  200
// @Router   /api/v1/deployments/{id}/events [get]
// @Success  200    {object}  api_response.Response{data=[]res.EventRes}
func (controller *EventController) EventIndex(context *gin.Context) {
	id := context.Param("id")
	page, _ := strconv.ParseInt(context.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(context.DefaultQuery("limit", "10"), 10, 64)
	if page < 1 {
		page = 1
	}

	events, pagination, err := controller.eventService.ListEvent(page, limit, id, context)
	if err != nil {
		errRes := api_response.BuildErrorResponse(
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			err.Error(),
			nil,
		)
		context.AbortWithStatusJSON(errRes.Code, errRes)
		return
	}
	eventListRes := mapper.ToDeploymentEventsList(events)
	apiRes := api_response.BuildPaginationResponse(
		http.StatusOK,
		http.StatusText(http.StatusOK),
		eventListRes,
		pagination,
	)

	context.JSON(apiRes.Code, apiRes)
}
