package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Event struct {
	Id             primitive.ObjectID `bson:"_id"`
	DeploymentId   primitive.ObjectID `bson:"deployment_id"`
	Type           string             `bson:"type"`
	TriggeredBy    string             `bson:"triggered_by"`
	TriggeredValue string             `bson:"triggered_value"`
	Reason         *string            `bson:"reason"`
	LatestStatus   string             `bson:"latest_status"`
	CreatedAt      time.Time          `bson:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at"`
}
