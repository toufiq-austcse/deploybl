package mapper

import (
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
	latestStatus string,
) *model.Event {
	return &model.Event{
		Type:           eventType,
		DeploymentId:   deploymentId,
		TriggeredBy:    triggeredBy,
		TriggeredValue: triggeredValue,
		Reason:         reason,
		LatestStatus:   latestStatus,
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
		LatestStatus:   event.LatestStatus,
		Reason:         event.Reason,
		CreatedAt:      event.CreatedAt,
		UpdatedAt:      event.UpdatedAt,
	}
}

func GenerateTitle(event *model.Event) string {
	switch event.Type {
	case deployment_events_enums.NEW_DEPLOYMENT:
		return "New deployment started"
	case deployment_events_enums.RESTARTED:
		return "Deployment restarted"
	case deployment_events_enums.STOPPED:
		return "Deployment stopped"
	case deployment_events_enums.INITIAL_DEPLPYMENT:
		return "Initial deployment started"
	default:
		return ""
	}
}
