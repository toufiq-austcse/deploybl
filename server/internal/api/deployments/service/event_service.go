package service

import (
	"context"
	"fmt"
	deployItConfig "github.com/toufiq-austcse/deployit/config"
	deployment_event_status_enums "github.com/toufiq-austcse/deployit/enums/deployment_event_status"
	"github.com/toufiq-austcse/deployit/pkg/aws/s3"
	"github.com/toufiq-austcse/deployit/pkg/utils"
	"math"
	"os"
	"time"

	"github.com/toufiq-austcse/deployit/pkg/api_response"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/toufiq-austcse/deployit/internal/api/deployments/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type EventService struct {
	eventCollection  *mongo.Collection
	dockerService    *DockerService
	s3ManagerService *s3.S3ManagerService
}

func NewEventService(
	database *mongo.Database,
	s3ManagerService *s3.S3ManagerService,
) *EventService {
	collection := database.Collection("events")
	return &EventService{eventCollection: collection, s3ManagerService: s3ManagerService}
}

func (service *EventService) Create(model *model.Event, ctx context.Context) error {
	currentTime := time.Now()
	model.CreatedAt = currentTime
	model.UpdatedAt = currentTime
	model.Id = primitive.NewObjectID()
	_, err := service.eventCollection.InsertOne(ctx, model)
	return err
}

func (service *EventService) FindLatestEventByDeploymentIdAndStatus(
	deploymentId primitive.ObjectID,
	status string,
	ctx context.Context,
) (*model.Event, error) {
	var event *model.Event

	opts := options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}})
	err := service.eventCollection.FindOne(
		ctx,
		bson.M{"deployment_id": deploymentId, "status": status},
		opts,
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
func (service *EventService) UpdateEvent(id primitive.ObjectID, updates map[string]interface{}, ctx context.Context) (*model.Event, error) {
	_, err := service.eventCollection.UpdateByID(
		ctx,
		id,
		bson.M{"$set": updates},
	)
	if err != nil {
		return nil, err
	}
	if updates["status"] == deployment_event_status_enums.SUCCESS || updates["status"] == deployment_event_status_enums.FAILED {
		go func() {
			fileKey, uploadErr := service.UploadLogFile(id.Hex())
			if uploadErr != nil {
				fmt.Println("error in uploading log file ", uploadErr.Error())
				return
			}
			_, updateErr := service.UpdateEvent(id, bson.M{"log_file_key": fileKey}, ctx)
			if updateErr != nil {
				fmt.Println("error in updating log file key ", updateErr.Error())
			}

		}()
	}
	return service.FindById(id)

}

func (service *EventService) UpdateStatusById(
	id primitive.ObjectID,
	status string,
	ctx context.Context,
) (*model.Event, error) {
	return service.UpdateEvent(id, bson.M{"status": status}, ctx)
}
func (service *EventService) UploadLogFile(eventId string) (*string, error) {
	logFilePath := utils.GetEventLogFilePath(eventId)
	fileKey, err := service.s3ManagerService.UploadFile(logFilePath, deployItConfig.AppConfig.AWS_CONFIG.AWS_S3_EVENT_LOG_PATH)
	if err != nil {
		return nil, err
	}
	return fileKey, nil
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

func (service *EventService) FindLatestProcessingEventsByDeploymentIds(
	deploymentIds []primitive.ObjectID,
	ctx context.Context,
) ([]*model.Event, error) {
	events := []*model.Event{}

	filter := bson.M{"deployment_id": bson.M{"$in": deploymentIds}, "status": "processing"}

	cursor, err := service.eventCollection.Find(ctx, filter)
	if err != nil {
		return events, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var event *model.Event
		decodeErr := cursor.Decode(&event)
		if decodeErr != nil {
			return events, decodeErr
		}
		events = append(events, event)
	}

	return events, nil
}

func (service *EventService) WriteToFile(text string, event *model.Event) {
	if event == nil {
		return
	}
	createErr := utils.CreateDirIfNotExists(deployItConfig.AppConfig.EVENT_LOGS_PATH)
	if createErr != nil {
		fmt.Println("error in creating dir ", createErr.Error())
		return
	}
	fileName := utils.GetEventLogFilePath(event.Id.Hex())
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Println("error in opening file ", err.Error())
	} else {
		defer file.Close()
		if _, writeErr := file.WriteString(text + "\n"); writeErr != nil {
			fmt.Println("error in writing file ", writeErr.Error())
		}
	}
}
