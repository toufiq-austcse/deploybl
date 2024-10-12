package model

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Id       primitive.ObjectID `bson:"_id"`
	Provider string             `bson:"provider"`
	Email    *string            `bson:"email"`
	Phone    *string            `bson:"phone"`
	Name     *string            `bson:"name"`
	UId      *string            `bson:"uid"`
	PhotoUrl *string            `bson:"photo_url"`
}

func CreateUserIndex(userCollection *mongo.Collection) {
	indexModel := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "uid", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	}
	_, err := userCollection.Indexes().CreateMany(context.Background(), indexModel)
	if err != nil {
		fmt.Println("error in user index create ", err.Error())
		return
	}
	fmt.Println("index created successfully in users")
}
