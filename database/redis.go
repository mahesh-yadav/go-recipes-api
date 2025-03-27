package database

import (
	"context"
	"log"

	"github.com/mahesh-yadav/go-recipes-api/config"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func ConnectToRedis(config *config.Config) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisUri,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	status := client.Ping(context.Background())
	redisClient = client

	log.Println("Successfully connected to Redis: ", status)
}

func GetRedisClient(config *config.Config) *redis.Client {
	if redisClient == nil {
		ConnectToRedis(config)
	}
	return redisClient
}
