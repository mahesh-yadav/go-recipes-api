package database

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/rs/zerolog/log"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/mahesh-yadav/go-recipes-api/config"
	"github.com/mahesh-yadav/go-recipes-api/models"
)

var mongoClient *mongo.Client

func ConnectToMongoDB(config *config.Config) {
	clientOptions := options.Client().SetTimeout(time.Duration(config.MongoServerSelectionTimeoutMS) * time.Millisecond).ApplyURI(config.MongoUri)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to MongoDB")
	}

	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		log.Fatal().Err(err).Msg("Error pinging to MongoDB")
	}

	mongoClient = client

	log.Info().Msg("Successfully connected to MongoDB!")
}

func GetMongoClient(config *config.Config) *mongo.Client {
	if mongoClient == nil {
		ConnectToMongoDB(config)
	}
	return mongoClient
}

func GetMongoCollection(config *config.Config, collectionName string) *mongo.Collection {
	client := GetMongoClient(config)
	return client.Database(config.MongoDBName).Collection(collectionName)
}

func InitDB(collection *mongo.Collection) {
	recipes := make([]models.Recipe, 0)
	file, _ := os.ReadFile("recipes.json")
	err := json.Unmarshal(file, &recipes)
	if err != nil {
		log.Fatal().Err(err).Msg("Error unmarshalling JSON file")
	}

	var listOfRecipes []interface{}
	for _, recipe := range recipes {
		listOfRecipes = append(listOfRecipes, recipe)
	}

	insertManyResult, err := collection.InsertMany(context.Background(), listOfRecipes)
	if err != nil {
		log.Fatal().Err(err).Msg("Error inserting documents to MongoDB")
	}

	log.Info().Int("Inserted recipes: ", len(insertManyResult.InsertedIDs))
}
