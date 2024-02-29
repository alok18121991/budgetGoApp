package configs

import (
	"context"
	"fmt"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient *mongo.Client
	once        sync.Once
)

// ConnectMongoDB connects to the MongoDB database
func ConnectMongoDB() *mongo.Client {
	once.Do(func() {
		clientOpts := options.Client().ApplyURI(GetMongoURI())
		client, err := mongo.Connect(context.Background(), clientOpts)
		if err != nil {
			log.Fatal(err)
		}

		// Ping the database to verify the connection
		err = client.Ping(context.Background(), nil)
		if err != nil {
			client.Disconnect(context.Background())
			log.Fatal(err)
		}

		fmt.Println("Connected to MongoDB")
		mongoClient = client
	})

	return mongoClient
}

// Client instance
var DB *mongo.Client = ConnectMongoDB()

// getting database collections
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("budget").Collection(collectionName)
	return collection
}

func GetSession(client *mongo.Client) (mongo.Session, error) {
	return client.StartSession()
}
