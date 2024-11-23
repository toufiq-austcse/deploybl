package controller

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"github.com/toufiq-austcse/deployit/internal/api/auth/dto/req"
	authMapper "github.com/toufiq-austcse/deployit/internal/api/auth/mapper"
	"github.com/toufiq-austcse/deployit/internal/api/users/mapper"
	userService "github.com/toufiq-austcse/deployit/internal/api/users/service"
	"github.com/toufiq-austcse/deployit/pkg/api_response"
	"github.com/toufiq-austcse/deployit/pkg/firebase"
	"net/http"
)

type AuthController struct {
	firebaseClient *firebase.Client
	userService    *userService.UserService
}

func NewAuthController(
	firebaseClient *firebase.Client, userService *userService.UserService) *AuthController {
	return &AuthController{
		firebaseClient: firebaseClient,
		userService:    userService,
	}
}

func (controller *AuthController) SignUp(context *gin.Context) {
	body := &req.SignUpReqDto{}
	if err := body.Validate(context); err != nil {
		errRes := api_response.BuildErrorResponse(
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
			err.Error(),
			nil,
		)
		context.AbortWithStatusJSON(http.StatusBadRequest, errRes)
		return
	}

	userToCreate := (&auth.UserToCreate{}).DisplayName(body.Name).Email(body.Email).Password(body.Password)
	createdUser, err := controller.firebaseClient.AuthClient.CreateUser(context, userToCreate)
	if err != nil {
		if err.Error() == "user with the provided email already exists" {
			errRes := api_response.BuildErrorResponse(
				http.StatusConflict,
				http.StatusText(http.StatusConflict),
				err.Error(),
				nil,
			)
			context.AbortWithStatusJSON(http.StatusConflict, errRes)
			return
		}
		errRes := api_response.BuildErrorResponse(
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			err.Error(),
			nil,
		)
		context.AbortWithStatusJSON(http.StatusInternalServerError, errRes)
		return
	}

	dbUserToCreate := mapper.MapFirebaseUserInfoToCreate(*createdUser.UserInfo)
	dbCreateErr := controller.userService.Create(dbUserToCreate, context)
	if dbCreateErr != nil {
		errRes := api_response.BuildErrorResponse(
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			dbCreateErr.Error(),
			nil,
		)
		context.AbortWithStatusJSON(http.StatusInternalServerError, errRes)
		return
	}

	token, tokenCreatedErr := controller.firebaseClient.AuthClient.CustomToken(context, createdUser.UID)
	if tokenCreatedErr != nil {
		errRes := api_response.BuildErrorResponse(
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			tokenCreatedErr.Error(),
			nil,
		)
		context.AbortWithStatusJSON(http.StatusInternalServerError, errRes)
		return
	}

	signUpRes := authMapper.ToSignUpResDto(token, *createdUser.UserInfo)
	apiRes := api_response.BuildResponse(
		http.StatusCreated,
		http.StatusText(http.StatusCreated),
		signUpRes,
	)

	context.JSON(apiRes.Code, apiRes)
}

func (controller *AuthController) Login(context *gin.Context) {

}
