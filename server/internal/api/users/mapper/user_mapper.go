package mapper

import (
	"firebase.google.com/go/v4/auth"
	"github.com/toufiq-austcse/deployit/internal/api/users/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MapFirebaseUserInfoToCreate(userInfo auth.UserInfo) *model.User {
	return &model.User{
		Id:       primitive.NewObjectID(),
		Provider: "firebase",
		Email:    &userInfo.Email,
		Name:     &userInfo.DisplayName,
		Phone:    &userInfo.PhoneNumber,
		PhotoUrl: &userInfo.PhotoURL,
		UId:      &userInfo.UID,
	}
}
