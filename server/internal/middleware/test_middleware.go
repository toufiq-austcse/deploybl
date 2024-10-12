package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func TestMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		fmt.Println("test middleware function called")
		context.Next()
	}
}
