package mapper

import (
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
