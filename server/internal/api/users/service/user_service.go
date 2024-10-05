package service

import (
	"context"
	"github.com/toufiq-austcse/deployit/internal/api/users/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	userCollection *mongo.Collection
}

func NewUserService(database *mongo.Database) *UserService {
	collection := database.Collection("users")
	go model.CreateUserIndex(collection)
	return &UserService{userCollection: collection}
}

func (service *UserService) Create(model *model.User, ctx context.Context) error {
	_, err := service.userCollection.InsertOne(ctx, model)
	return err
}

func (service *UserService) FindUserByUId(uId string, ctx context.Context) *model.User {
	var user *model.User
	filter := bson.M{"uid": uId}
	err := service.userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil
	}
	return user
}
