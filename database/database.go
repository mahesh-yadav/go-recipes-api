package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/mahesh-yadav/go-recipes-api/config"
)

var mongoClient *mongo.Client

func ConnectToDB(config *config.Config) {
	clientOptions := options.Client().SetTimeout(10 * time.Second).ApplyURI(config.MongoUri)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal("Error connecting to MongoDB", err)
	}

	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		panic(err)
	}

	mongoClient = client

	log.Println("Successfully connected to MongoDB!")
}

func GetMongoClient(config *config.Config) *mongo.Client {
	if mongoClient == nil {
		ConnectToDB(config)
	}
	return mongoClient
}

func GetCollection(config *config.Config, collectionName string) *mongo.Collection {
	client := GetMongoClient(config)
	return client.Database(config.MongoDBName).Collection(collectionName)
}
