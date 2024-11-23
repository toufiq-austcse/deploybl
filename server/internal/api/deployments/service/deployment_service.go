package service

import (
	"context"
	"fmt"
	"math"
	"time"

	deployment_event_status_enums "github.com/toufiq-austcse/deployit/enums/deployment_event_status"

	deployment_events_enums "github.com/toufiq-austcse/deployit/enums/deployment_events"
	"github.com/toufiq-austcse/deployit/enums/deployment_events_triggered_by"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/mapper"

	"github.com/toufiq-austcse/deployit/enums"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/model"
	"github.com/toufiq-austcse/deployit/pkg/api_response"
	"github.com/toufiq-austcse/deployit/pkg/app_errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DeploymentService struct {
	deploymentCollection *mongo.Collection
	dockerService        *DockerService
	eventService         *EventService
}

func NewDeploymentService(
	database *mongo.Database,
	dockerService *DockerService,
	eventService *EventService,
) *DeploymentService {
	collection := database.Collection("deployments")
	go model.CreateDeploymentIndex(collection)
	return &DeploymentService{
		deploymentCollection: collection,
		dockerService:        dockerService,
		eventService:         eventService,
	}
}

func (service *DeploymentService) Create(
	model *model.Deployment,
	ctx context.Context,
) (*model.Event, error) {
	model.Id = primitive.NewObjectID()
	currentTime := time.Now()
	model.CreatedAt = currentTime
	model.UpdatedAt = currentTime
	model.LastDeploymentInitiatedAt = &currentTime
	_, err := service.deploymentCollection.InsertOne(ctx, model)
	if err != nil {
		return nil, err
	}
	newEventModel := mapper.MapEventModelToSave(
		model.Id,
		deployment_events_enums.INITIAL_DEPLPYMENT,
		deployment_events_triggered_by.USER,
		model.UserId.Hex(),
		nil,
	)
	_ = service.eventService.Create(newEventModel, ctx)

	return newEventModel, err
}

func (service *DeploymentService) FindBySubDomainName(
	domainName *string,
	ctx context.Context,
) *model.Deployment {
	var deployment *model.Deployment
	filter := bson.M{"sub_domain_name": domainName}
	err := service.deploymentCollection.FindOne(ctx, filter).Decode(&deployment)
	if err != nil {
		return nil
	}
	return deployment
}

func (service *DeploymentService) FindById(id string, ctx context.Context) *model.Deployment {
	var deployment *model.Deployment
	oId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil
	}

	filter := bson.M{"_id": oId}
	err = service.deploymentCollection.FindOne(ctx, filter).Decode(&deployment)
	if err != nil {
		return nil
	}
	return deployment
}

func (service *DeploymentService) FindUserDeploymentById(
	id string,
	userId primitive.ObjectID,
	ctx context.Context,
) *model.Deployment {
	var deployment *model.Deployment
	oId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil
	}

	filter := bson.M{"_id": oId, "user_id": userId}
	err = service.deploymentCollection.FindOne(ctx, filter).Decode(&deployment)
	if err != nil {
		return nil
	}
	return deployment
}

func (service *DeploymentService) ListDeployment(
	page, limit int64,
	userId primitive.ObjectID,
	ctx context.Context,
) ([]*model.Deployment, *api_response.Pagination, error) {
	deployments := []*model.Deployment{}

	filter := bson.M{"user_id": userId}

	totalDocs, err := service.deploymentCollection.CountDocuments(ctx, filter)
	if err != nil {
		return deployments, nil, err
	}
	lastPage := int64(math.Ceil(float64(totalDocs) / float64(limit)))

	findOptions := options.Find()
	skip := page*limit - limit
	findOptions.SetSort(bson.M{"created_at": -1}).SetLimit(limit).SetSkip(skip)

	cursor, err := service.deploymentCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return deployments, nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var deployment *model.Deployment
		decodeErr := cursor.Decode(&deployment)
		if decodeErr != nil {
			return deployments, nil, decodeErr
		}
		deployments = append(deployments, deployment)
	}

	return deployments, &api_response.Pagination{
		Total:       totalDocs,
		CurrentPage: page,
		LastPage:    lastPage,
		PerPage:     limit,
	}, nil
}

func (service *DeploymentService) UpdateDeployment(
	deploymentId string,
	updates map[string]interface{},
	existingEvent *model.Event,
	ctx context.Context,
) (*model.Deployment, error) {
	updates["updated_at"] = time.Now()
	if updates["latest_status"] == enums.QUEUED {
		updates["last_deployment_initiated_at"] = time.Now()
	}
	fmt.Println("updating ", deploymentId, updates)
	oId, err := primitive.ObjectIDFromHex(deploymentId)
	if err != nil {
		return nil, err
	}
	updatedResult, err := service.deploymentCollection.UpdateByID(ctx, oId, bson.M{
		"$set": updates,
	})
	if err != nil {
		return nil, err
	}
	if updatedResult.MatchedCount != 1 {
		return nil, app_errors.CannotUpdateError
	}
	updatedDeployment := service.FindById(deploymentId, ctx)

	if updates["latest_status"] != nil && existingEvent != nil {
		eventStatus := service.GetEventStatusByDeploymentStatus(existingEvent, updatedDeployment.LatestStatus)
		fmt.Println("event status ", eventStatus)
		service.eventService.UpdateStatusById(existingEvent.Id, eventStatus, ctx)
	}
	return updatedDeployment, nil
}

func (service *DeploymentService) GetLatestStatusByIds(
	ids []string,
	userId primitive.ObjectID,
	ctx context.Context,
) ([]*model.Deployment, error) {
	objectIds := []primitive.ObjectID{}
	deployments := []*model.Deployment{}

	for _, id := range ids {
		oId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return deployments, err
		}
		objectIds = append(objectIds, oId)
	}
	filter := bson.M{"_id": bson.M{"$in": objectIds}, "user_id": userId}
	projection := bson.M{
		"latest_status":    1,
		"last_deployed_at": 1,
		"_id":              1,
		"sub_domain_name":  1,
	}
	findOptions := options.Find().SetProjection(projection)

	cursor, err := service.deploymentCollection.Find(ctx, filter, findOptions)
	defer cursor.Close(ctx)
	if err != nil {
		return deployments, err
	}

	for cursor.Next(ctx) {
		var deployment *model.Deployment
		decodeErr := cursor.Decode(&deployment)
		if decodeErr != nil {
			return deployments, decodeErr
		}
		deployments = append(deployments, deployment)
	}

	return deployments, nil
}

func (service *DeploymentService) FindByDeploymentStatus(status string) []model.Deployment {
	var deployments []model.Deployment
	filter := bson.M{"latest_status": status}
	cursor, err := service.deploymentCollection.Find(context.Background(), filter)
	if err != nil {
		return deployments
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var deployment model.Deployment
		if decodeErr := cursor.Decode(&deployment); decodeErr != nil {
			continue
		}
		deployments = append(deployments, deployment)
	}
	return deployments
}

func (service *DeploymentService) GetContainerIdsFromDeployments(
	deployments []model.Deployment,
) []string {
	var containerIds []string
	for _, deployment := range deployments {
		if deployment.ContainerId != nil {
			containerIds = append(containerIds, *deployment.ContainerId)
		}
	}
	return containerIds
}

func (service *DeploymentService) UpdateDeploymentStatusByContainerIds(
	skipContainerIds []string,
	currentStatus string,
	updatedStatus string,
	ctx context.Context,
) (int64, error) {
	filter := bson.M{
		"container_id":  bson.M{"$nin": skipContainerIds},
		"latest_status": currentStatus,
	}
	update := bson.M{"$set": bson.M{"latest_status": updatedStatus, "updated_at": time.Now()}}
	result, err := service.deploymentCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

func (service *DeploymentService) UpdateLatestStatus(
	deploymentId string,
	status string,
	existingEvent *model.Event,
	context context.Context,
) (*model.Deployment, error) {
	return service.UpdateDeployment(deploymentId, map[string]interface{}{
		"latest_status": status,
	}, existingEvent, context)
}

func (service *DeploymentService) IsRestartable(deployment *model.Deployment) bool {
	if deployment.LatestStatus == enums.LIVE || deployment.LatestStatus == enums.STOPPED ||
		deployment.LatestStatus == enums.FAILED {
		return true
	}
	return false
}

func (service *DeploymentService) IsRebuildAble(deployment *model.Deployment) bool {
	if deployment.LatestStatus == enums.LIVE || deployment.LatestStatus == enums.STOPPED ||
		deployment.LatestStatus == enums.FAILED {
		return true
	}
	return false
}

func (service *DeploymentService) IsStopAble(deployment *model.Deployment) bool {
	if deployment.LatestStatus == enums.LIVE {
		return true
	}
	return false
}

func (service *DeploymentService) CountDeploymentByRepositoryName(
	repositoryName string,
	ctx context.Context,
) (int64, error) {
	filter := bson.M{"repository_name": repositoryName}
	return service.deploymentCollection.CountDocuments(ctx, filter)
}

func (service *DeploymentService) GetDeploymentsByIds(
	ids []primitive.ObjectID,
	ctx context.Context,
) ([]model.Deployment, error) {
	var deployments []model.Deployment
	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := service.deploymentCollection.Find(context.Background(), filter)
	if err != nil {
		return deployments, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var deployment model.Deployment
		if decodeErr := cursor.Decode(&deployment); decodeErr != nil {
			continue
		}
		deployments = append(deployments, deployment)
	}
	return deployments, nil
}

func (service *DeploymentService) GetEventStatusByDeploymentStatus(
	event *model.Event,
	deploymentStatus string,
) string {
	fmt.Println("event type ", event.Type, " deployment status ", deploymentStatus)
	switch event.Type {
	case deployment_events_enums.RESTART_DEPLOYMENT:
		if deploymentStatus == enums.LIVE {
			return deployment_event_status_enums.SUCCESS
		} else if deploymentStatus == enums.FAILED {
			return deployment_event_status_enums.FAILED
		} else {
			return deployment_event_status_enums.PROCESSING
		}
	case deployment_events_enums.REBUILD_DEPLOYMENT:
		if deploymentStatus == enums.LIVE {
			return deployment_event_status_enums.SUCCESS
		} else if deploymentStatus == enums.FAILED {
			return deployment_event_status_enums.FAILED
		} else {
			return deployment_event_status_enums.PROCESSING
		}
	case deployment_events_enums.INITIAL_DEPLPYMENT:
		if deploymentStatus == enums.LIVE {
			return deployment_event_status_enums.SUCCESS
		} else if deploymentStatus == enums.FAILED {
			return deployment_event_status_enums.FAILED
		} else {
			return deployment_event_status_enums.PROCESSING
		}
	case deployment_events_enums.STOP_DEPLOYMENT:
		if deploymentStatus == enums.STOPPED {
			return deployment_event_status_enums.SUCCESS
		} else if deploymentStatus == enums.FAILED {
			return deployment_event_status_enums.FAILED
		} else {
			return deployment_event_status_enums.PROCESSING
		}
	default:
		return ""

	}
}
