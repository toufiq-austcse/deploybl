package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/toufiq-austcse/go-api-boilerplate/ent"
	"github.com/toufiq-austcse/go-api-boilerplate/internal/api/todo/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	client *mongo.Client
}

func NewRepository(client *mongo.Client) *Repository {
	return &Repository{
		client: client,
	}
}
func (repository *Repository) Create(model *models.CreateTodoModel, context context.Context) (*ent.Todo, error) {
	return nil, nil
}

func (repository *Repository) GetAll(ctx context.Context) ([]*ent.Todo, error) {
	return []*ent.Todo{}, nil
}

func (repository *Repository) FindOne(id uuid.UUID, ctx context.Context) (*ent.Todo, error) {
	return nil, nil
}
