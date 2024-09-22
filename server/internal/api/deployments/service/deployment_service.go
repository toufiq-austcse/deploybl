package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/model"
	"github.com/toufiq-austcse/deployit/pkg/api_response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DeploymentService struct {
	deploymentCollection *mongo.Collection
}

func NewDeploymentService(database *mongo.Database) *DeploymentService {
	collection := database.Collection("deployments")
	model.CreateDeploymentIndex(collection)
	return &DeploymentService{deploymentCollection: collection}
}

func (service *DeploymentService) Create(model *model.Deployment, ctx context.Context) error {
	_, err := service.deploymentCollection.InsertOne(ctx, model)
	return err
}

func (service *DeploymentService) FindBySubDomainName(domainName *string, ctx context.Context) *model.Deployment {
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

func (service *DeploymentService) ListDeployment(page, limit int64, ctx context.Context) ([]*model.Deployment, *api_response.Pagination, error) {
	deployments := []*model.Deployment{}

	totalDocs, err := service.deploymentCollection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return deployments, nil, err
	}
	lastPage := int64(float64(totalDocs) / float64(totalDocs))

	findOptions := options.Find()
	skip := int64(page*limit - limit)
	findOptions.SetSort(bson.M{"created_at": -1}).SetLimit(limit).SetSkip(skip)

	cursor, err := service.deploymentCollection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return deployments, nil, err
	}
	for cursor.Next(ctx) {
		var deployment *model.Deployment
		decodeErr := cursor.Decode(&deployment)
		if decodeErr != nil {
			return deployments, nil, decodeErr

		}
		deployments = append(deployments, deployment)
	}
	cursor.Close(ctx)

	return deployments, &api_response.Pagination{
		Total:       totalDocs,
		CurrentPage: page,
		LastPage:    lastPage,
		PerPage:     limit,
	}, nil
}

func (service *DeploymentService) UpdateDeployment(deploymentId string, updates map[string]interface{}, ctx context.Context) (*model.Deployment, error) {
	fmt.Println("updating ", deploymentId, updates)
	oId, err := primitive.ObjectIDFromHex(deploymentId)
	if err != nil {
		return nil, err
	}
	updatedResult, err := service.deploymentCollection.UpdateByID(ctx, oId, bson.M{
		"$set": updates,
	})
	if err != nil {
		fmt.Println("err in updating deployment", err)
		return nil, err
	}
	if updatedResult.MatchedCount != 1 {
		return nil, errors.New("update error")
	}
	return service.FindById(deploymentId, ctx), err
}
