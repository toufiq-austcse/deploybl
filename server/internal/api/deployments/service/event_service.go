package service

import (
	"context"
	"github.com/toufiq-austcse/deployit/pkg/api_response"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"

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

func (service *EventService) FindById(id primitive.ObjectID) (*model.Event, error) {
	var event *model.Event
	filter := bson.M{"_id": id}
	err := service.eventCollection.FindOne(context.Background(), filter).Decode(&event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (service *EventService) UpdateLatestStatusByDeploymentId(
	deploymentId primitive.ObjectID,
	latestStatus string,
	ctx context.Context,
) (*model.Event, error) {
	_, err := service.eventCollection.UpdateOne(
		ctx,
		bson.M{"deployment_id": deploymentId},
		bson.M{"$set": bson.M{"latest_status": latestStatus}},
	)
	if err != nil {
		return nil, err
	}
	return service.FindById(deploymentId)
}

func (service *EventService) ListEvent(
	page, limit int64,
	deploymentId string,
	ctx context.Context,
) ([]*model.Event, *api_response.Pagination, error) {
	deploymentIdObj, err := primitive.ObjectIDFromHex(deploymentId)
	if err != nil {
		return nil, nil, err
	}
	events := []*model.Event{}

	filter := bson.M{"deployment_id": deploymentIdObj}

	totalDocs, err := service.eventCollection.CountDocuments(ctx, filter)
	if err != nil {
		return events, nil, err
	}
	lastPage := int64(math.Ceil(float64(totalDocs) / float64(limit)))

	findOptions := options.Find()
	skip := page*limit - limit
	findOptions.SetSort(bson.M{"created_at": -1}).SetLimit(limit).SetSkip(skip)

	cursor, err := service.eventCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return events, nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var event *model.Event
		decodeErr := cursor.Decode(&event)
		if decodeErr != nil {
			return events, nil, decodeErr
		}
		events = append(events, event)
	}

	return events, &api_response.Pagination{
		Total:       totalDocs,
		CurrentPage: page,
		LastPage:    lastPage,
		PerPage:     limit,
	}, nil
}
