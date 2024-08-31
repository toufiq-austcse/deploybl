package res

import "github.com/google/uuid"

type TodoRespData struct {
	Id     uuid.UUID `json:"id"`
	Title  string    `json:"title"`
	Status string    `json:"status"`
}
