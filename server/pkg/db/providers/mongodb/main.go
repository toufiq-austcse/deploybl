package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/toufiq-austcse/deployit/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func New() (*mongo.Client, *mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig.MONGO_DB_CONFIG.URL))
	if err != nil {
		return nil, nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("Connected to mongodb...")

	return client, client.Database(config.AppConfig.MONGO_DB_CONFIG.DB_NAME), nil
}
