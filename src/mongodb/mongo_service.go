package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var Client *mongo.Client
var Database *mongo.Database

func Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI(fmt.Sprintf("mongodb://%s:%s", env.dns, env.port)).
		SetAuth(options.Credential{Username: env.username, Password: env.password}))
	if err != nil {
		return err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}
	Client = client
	Database = client.Database("lambda-clone")
	return nil
}

func Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return Client.Disconnect(ctx)
}
