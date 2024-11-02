package mapper

import (
	deployment_event_status_enums "github.com/toufiq-austcse/deployit/enums/deployment_event_status"
	deployment_events_enums "github.com/toufiq-austcse/deployit/enums/deployment_events"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/dto/res"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MapEventModelToSave(
	deploymentId primitive.ObjectID,
	eventType string,
	triggeredBy string,
	triggeredValue string,
	reason *string,
) *model.Event {
	return &model.Event{
		Type:           eventType,
		DeploymentId:   deploymentId,
		TriggeredBy:    triggeredBy,
		TriggeredValue: triggeredValue,
		Reason:         reason,
		Status:         deployment_event_status_enums.PROCESSING,
	}
}

func ToDeploymentEventsList(events []*model.Event) []res.EventRes {
	var eventResList []res.EventRes = []res.EventRes{}
	for _, event := range events {
		eventResList = append(eventResList, ToDeploymentEvent(event))
	}
	return eventResList
}

func ToDeploymentEvent(event *model.Event) res.EventRes {
	return res.EventRes{
		Id:             event.Id.Hex(),
		DeploymentId:   event.DeploymentId.Hex(),
		Title:          GenerateTitle(event),
		Type:           event.Type,
		TriggeredBy:    event.TriggeredBy,
		TriggeredValue: event.TriggeredValue,
		Status:         event.Status,
		Reason:         event.Reason,
		CreatedAt:      event.CreatedAt,
		UpdatedAt:      event.UpdatedAt,
	}
}

func GenerateTitle(event *model.Event) string {
	switch event.Type {
	case deployment_events_enums.RESTART_DEPLOYMENT:
		return "Restart deployment"
	case deployment_events_enums.STOP_DEPLOYMENT:
		return "Stop deployment"
	case deployment_events_enums.INITIAL_DEPLPYMENT:
		return "First deployment started"
	case deployment_events_enums.REBUILD_DEPLOYMENT:
		return "Rebuild deployment"
	default:
		return ""
	}
}
