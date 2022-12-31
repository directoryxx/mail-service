package infrastructure

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo() (client *mongo.Client, err error) {
	dsn := "mongodb://" + os.Getenv("DB_USERNAME") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/admin"

	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(dsn))
	if err != nil {
		return nil, err
	}

	return client, nil
}
