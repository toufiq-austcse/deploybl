package payloads

import "go.mongodb.org/mongo-driver/bson/primitive"

type PreRunRepoWorkerPayload struct {
	DeploymentId string             `json:"deployment_id"`
	EventId      primitive.ObjectID `json:"event_id"`
}
