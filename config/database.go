package config

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type Database struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func ConnectDB() (*Database, error) {
	// Load the MongoDB URI from the environment variable
	mongoURI := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("MONGODB_DBNAME")
	fmt.Print("Connecting to MongoDB... ")
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}

	db := client.Database(dbName)
	return &Database{
		Client: client,
		DB:     db}, nil

}

func (d *Database) Ping() error {
	return d.Client.Ping(context.TODO(), readpref.Primary())
}

func (d *Database) Close() error {
	if err := d.Client.Disconnect(context.TODO()); err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %v", err)
	}
	return nil
}

func (d *Database) GetCollection(collectionName string) *mongo.Collection {
	return d.DB.Collection(collectionName)
}
