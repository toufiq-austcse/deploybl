package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/toufiq-austcse/go-api-boilerplate/ent"
	"github.com/toufiq-austcse/go-api-boilerplate/internal/api/todo/apimodels/req"
	"github.com/toufiq-austcse/go-api-boilerplate/internal/api/todo/apimodels/res"
	"github.com/toufiq-austcse/go-api-boilerplate/internal/api/todo/models"
	"github.com/toufiq-austcse/go-api-boilerplate/internal/api/todo/repository"
)

type TodoRepository interface {
	Create(model *models.CreateTodoModel, context context.Context) (*ent.Todo, error)
	GetAll(ctx context.Context) ([]*ent.Todo, error)
	FindOne(uuid uuid.UUID, ctx context.Context) (*ent.Todo, error)
}
type TodoService struct {
	TodoRepo TodoRepository
}

func NewTodoService(repository *repository.Repository) *TodoService {
	return &TodoService{
		TodoRepo: repository,
	}
}

func (todoService *TodoService) TodoCreate(reqModel *req.TodoCreateReqModel, ctx context.Context) (res.TodoRespData, error) {
	createdTodo, err := todoService.TodoRepo.Create(&models.CreateTodoModel{
		Title:  reqModel.Title,
		Status: reqModel.Status,
	}, ctx)
	if err != nil {
		return res.TodoRespData{}, err
	}
	return res.TodoRespData{
		Id:     createdTodo.ID,
		Title:  createdTodo.Title,
		Status: createdTodo.Status,
	}, nil
}

func (todoService *TodoService) TodoList(ctx context.Context) ([]res.TodoRespData, error) {
	var todoRes []res.TodoRespData
	allTodos, err := todoService.TodoRepo.GetAll(ctx)
	if err != nil {
		return todoRes, err
	}
	for _, todo := range allTodos {
		todoRes = append(todoRes, res.TodoRespData{
			Title:  todo.Title,
			Status: todo.Status,
			Id:     todo.ID,
		})
	}
	fmt.Println(len(todoRes))
	return todoRes, nil
}

func (todoService *TodoService) TodoGet(uuid uuid.UUID, ctx context.Context) (*ent.Todo, error) {
	todo, err := todoService.TodoRepo.FindOne(uuid, ctx)
	if err != nil {
		if err.Error() == "ent: todo not found" {
			return nil, nil
		}
		return nil, err
	}
	return todo, err
}
