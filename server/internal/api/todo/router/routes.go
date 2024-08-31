package router

import (
	"github.com/gin-gonic/gin"
	"github.com/toufiq-austcse/go-api-boilerplate/internal/api/todo/controller"
	"github.com/toufiq-austcse/go-api-boilerplate/internal/middleware"
)

func Setup(group *gin.RouterGroup, controller *controller.TodoController) {
	group.Use(middleware.TestMiddleware())
	group.POST("todos", controller.CreateTodo)
	group.GET("todos", controller.ListTodo)
	group.GET("todos/:id", controller.GetOneTodo)
}
