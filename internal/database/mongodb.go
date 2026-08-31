package database

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoStorage struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func NewConnection() (*MongoStorage, error) {
	mongoURL := os.Getenv("MONGO_URL")
	if mongoURL == "" {
		return nil, fmt.Errorf("MONGO_URL environment variable is not set")
	}

	dbName := os.Getenv("MONGO_DB_NAME")
	if dbName == "" {
		return nil, fmt.Errorf("MONGO_DB_NAME is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Client().ApplyURI(mongoURL)
	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create mongo client: %w", err)
	}

	var result bson.M
	if err := client.Database("admin").RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		_ = client.Disconnect(ctx)
		return nil, fmt.Errorf("failed to ping mongodb: %w", err)
	}

	slog.Info("Successfully connected and pinged MongoDB")

	return &MongoStorage{
		Client:   client,
		Database: client.Database(dbName),
	}, nil
}
