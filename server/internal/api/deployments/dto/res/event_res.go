package res

import "time"

type EventRes struct {
	Id             string    `json:"id"`
	DeploymentId   string    `json:"deployment_id"`
	Title          string    `json:"title"`
	Type           string    `json:"type"`
	TriggeredBy    string    `json:"triggered_by"`
	TriggeredValue string    `json:"triggered_value"`
	LatestStatus   string    `json:"latest_status"`
	Reason         *string   `json:"reason"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}