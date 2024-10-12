package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/toufiq-austcse/deployit/internal/api/users/mapper"
	"github.com/toufiq-austcse/deployit/internal/api/users/service"
	"github.com/toufiq-austcse/deployit/pkg/api_response"
	"github.com/toufiq-austcse/deployit/pkg/firebase"
)

func AuthMiddleware(firebaseClient *firebase.Client, userService *service.UserService) gin.HandlerFunc {
	return func(context *gin.Context) {
		authToken := context.GetHeader("Authorization")
		if authToken == "" {
			apiRes := api_response.BuildErrorResponse(http.StatusUnauthorized, "Authorization token required", "", nil)
			context.AbortWithStatusJSON(apiRes.Code, apiRes)
			return
		}
		authToken = strings.TrimPrefix(authToken, "Bearer ")

		verifyRes, err := firebaseClient.AuthClient.VerifyIDToken(context, authToken)
		if err != nil {
			apiRes := api_response.BuildErrorResponse(http.StatusUnauthorized, "Invalid token", err.Error(), nil)
			context.AbortWithStatusJSON(apiRes.Code, apiRes)
			return

		}
		user := userService.FindUserByUId(verifyRes.UID, context)
		if user == nil {
			record, getUserErr := firebaseClient.AuthClient.GetUser(context, verifyRes.UID)
			if getUserErr != nil {
				apiRes := api_response.BuildErrorResponse(http.StatusUnauthorized, "Invalid token", err.Error(), nil)
				context.AbortWithStatusJSON(apiRes.Code, apiRes)
				return
			}
			newUser := mapper.MapFirebaseUserInfoToCreate(*record.UserInfo)
			createErr := userService.Create(newUser, context)
			if createErr != nil {
				apiRes := api_response.BuildErrorResponse(
					http.StatusUnauthorized,
					http.StatusText(http.StatusUnauthorized),
					createErr.Error(),
					nil,
				)
				context.AbortWithStatusJSON(apiRes.Code, apiRes)
				return
			}
			user = newUser
		}

		context.Set("user", user)
		context.Next()
	}
}
