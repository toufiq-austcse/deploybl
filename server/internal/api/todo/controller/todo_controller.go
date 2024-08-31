package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/toufiq-austcse/go-api-boilerplate/internal/api/todo/apimodels/req"
	"github.com/toufiq-austcse/go-api-boilerplate/internal/api/todo/service"
	"github.com/toufiq-austcse/go-api-boilerplate/pkg/api_response"
	"net/http"
)

type TodoController struct {
	TodoService *service.TodoService
}

func NewTodoController(todoService *service.TodoService) *TodoController {
	return &TodoController{
		TodoService: todoService,
	}
}

// CreateTodo hosts godoc
// @Summary  Create Todo
// @Tags     Todo
// @Accept   json
// @Produce  json
// @Param    request  body      req.TodoCreateReqModel  true  "Create Todo Payload"
// @Success  201      {object}  api_response.Response{data=res.TodoRespData}
// @Router   /api/v1/todos [post]
func (controller *TodoController) CreateTodo(context *gin.Context) {
	body := &req.TodoCreateReqModel{}
	if err := body.Validate(context); err != nil {
		res := api_response.BuildErrorResponse(http.StatusBadRequest, "Bad Request", err.Error(), nil)
		context.JSON(http.StatusBadRequest, res)
		return
	}
	create, err := controller.TodoService.TodoCreate(body, context)
	if err != nil {
		res := api_response.BuildErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error(), nil)
		context.JSON(http.StatusInternalServerError, res)
		return
	}
	res := api_response.BuildResponse(http.StatusCreated, "Todo Created", create)
	context.JSON(http.StatusCreated, res)
}

// ListTodo hosts godoc
// @Summary  List Todo
// @Tags     Todo
// @Accept   json
// @Produce  json
// @Success  200  {object}  api_response.Response{data=[]res.TodoRespData}
// @Router   /api/v1/todos [get]
func (controller *TodoController) ListTodo(context *gin.Context) {
	todos, err := controller.TodoService.TodoList(context)
	if err != nil {
		res := api_response.BuildErrorResponse(http.StatusInternalServerError, "Internal Server Error", err.Error(), nil)
		context.JSON(http.StatusInternalServerError, res)
		return
	}
	res := api_response.BuildResponse(http.StatusOK, "Todo List", todos)
	context.JSON(http.StatusOK, res)
}

// GetOneTodo hosts godoc
// @Summary  List Todo
// @Tags     Todo
// @Accept   json
// @Produce  json
// @Param    id   path      string  true  "todo id"
// @Success  200  {object}  api_response.Response{data=res.TodoRespData}
// @Router   /api/v1/todos/{id} [get]
func (controller *TodoController) GetOneTodo(context *gin.Context) {
	id := context.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		res := api_response.BuildErrorResponse(http.StatusBadRequest, "", err.Error(), nil)
		context.JSON(http.StatusBadRequest, res)
		return
	}

	todo, err := controller.TodoService.TodoGet(uuid, context)
	if err != nil {
		res := api_response.BuildErrorResponse(http.StatusInternalServerError, "", err.Error(), nil)
		context.JSON(http.StatusInternalServerError, res)
		return
	}
	if todo == nil {
		res := api_response.BuildErrorResponse(http.StatusNotFound, "", "Todo Not Found", nil)
		context.JSON(http.StatusNotFound, res)
		return
	}
	res := api_response.BuildResponse(http.StatusOK, "", todo)
	context.JSON(http.StatusOK, res)
}
