package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/toufiq-austcse/deployit/internal/api/deployments/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type EventService struct {
	eventCollection *mongo.Collection
	dockerService   *DockerService
}

func NewEventService(
	database *mongo.Database,
) *EventService {
	collection := database.Collection("events")
	return &EventService{eventCollection: collection}
}

func (service *EventService) Create(model *model.Event, ctx context.Context) error {
	currentTime := time.Now()
	model.CreatedAt = currentTime
	model.UpdatedAt = currentTime
	model.Id = primitive.NewObjectID()
	_, err := service.eventCollection.InsertOne(ctx, model)
	return err
}
func (service *EventService) FindByDeploymentIdAndLatestStatus(
	deploymentId primitive.ObjectID,
	latestStatus string,
	ctx context.Context,
) (*model.Event, error) {
	var event *model.Event
	err := service.eventCollection.FindOne(
		ctx,
		model.Event{DeploymentId: deploymentId, LatestStatus: latestStatus},
	).Decode(&event)
	if err != nil {
		return nil, err
	}
	return event, nil
}
