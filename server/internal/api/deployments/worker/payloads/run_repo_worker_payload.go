package payloads

import "go.mongodb.org/mongo-driver/bson/primitive"

type RunRepoWorkerPayload struct {
	DeploymentId string             `json:"deployment_id"`
	EventId      primitive.ObjectID `json:"event_id"`
}
