package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Database holds the MongoDB client and database reference.
type Database struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// NewDatabase creates a new MongoDB connection and returns a Database instance.
func NewDatabase(ctx context.Context, mongoURI string, dbName string) (*Database, error) {
	// create MongoDB client
	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// ping the database to verify the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	// get the database reference
	database := client.Database(dbName)

	return &Database{
		Client:   client,
		Database: database,
	}, nil
}

// Close closes the MongoDB connection.
func (db *Database) Close(ctx context.Context) error {
	if db.Client != nil {
		return db.Client.Disconnect(ctx)
	}
	return nil
}

// GetCollection returns a collection from the database.
func (db *Database) GetCollection(collectionName string) *mongo.Collection {
	return db.Database.Collection(collectionName)
}

