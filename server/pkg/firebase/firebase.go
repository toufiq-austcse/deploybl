package firebase

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"google.golang.org/api/option"
)

type Client struct {
	AuthClient *auth.Client
}

func NewFirebaseClient() (*Client, error) {
	opt := option.WithCredentialsFile("deploybl-7a03d-firebase-adminsdk-snvhm-7e1ee4d952.json")
	newApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		fmt.Println("error in creating firebase app ", err.Error())
		return nil, err
	}
	authClient, err := newApp.Auth(context.Background())
	if err != nil {
		fmt.Println("error in getting firebase auth ", err.Error())
		return nil, err
	}
	return &Client{
		AuthClient: authClient,
	}, nil
}
