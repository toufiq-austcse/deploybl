package service

import (
	"context"
	"github.com/toufiq-austcse/deployit/internal/api/deployments/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
