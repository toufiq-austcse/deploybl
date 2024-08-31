package model

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Deployment struct {
	Id                 primitive.ObjectID      `bson:"_id"`
	Title              string                  `bson:"title"`
	SubDomainName      string                  `bson:"sub_domain_name"`
	LatestStatus       string                  `bson:"latest_status"`
	LastDeployedAt     *time.Time              `bson:"last_deployed_at"`
	RepositoryProvider string                  `bson:"repository_provider"`
	RepositoryUrl      string                  `bson:"repository_url"`
	BranchName         string                  `bson:"branch_name"`
	DockerFilePath     string                  `bson:"docker_file_path"`
	DockerImageTag     *string                 `bson:"docker_image_tag"`
	ContainerId        *string                 `bson:"container_id"`
	Env                *map[string]interface{} `bson:"env"`
	CreatedAt          time.Time               `bson:"created_at"`
	UpdatedAt          time.Time               `bson:"updated_at"`
}

func CreateDeploymentIndex(deploymentCollection *mongo.Collection) {
	indexModel := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "sub_domain_name", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	}
	_, err := deploymentCollection.Indexes().CreateMany(context.Background(), indexModel)
	if err != nil {
		fmt.Println("error in deployment index create ", err.Error())
		return
	}
	fmt.Println("index created successfully in deployments")

}
