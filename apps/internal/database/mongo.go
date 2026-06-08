package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// MongoClient holds the global MongoDB client for health checks and shutdown.
var MongoClient *mongo.Client

// NewMongoClient connects to MongoDB and returns the client and database handle.
func NewMongoClient(mongoURI, dbName string) (*mongo.Client, *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().
		ApplyURI(mongoURI).
		SetMaxPoolSize(25).
		SetMinPoolSize(5).
		SetMaxConnIdleTime(5 * time.Minute)

	client, err := mongo.Connect(clientOpts)
	if err != nil {
		log.Fatalf("failed to connect to mongodb: %v", err)
	}

	// Verify connectivity
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("failed to ping mongodb: %v", err)
	}

	MongoClient = client

	return client, client.Database(dbName)
}
