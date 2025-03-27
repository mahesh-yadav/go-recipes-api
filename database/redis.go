package database

import (
	"context"
	"log"

	"github.com/mahesh-yadav/go-recipes-api/config"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func ConnectToRedis(config *config.Config) {
	if config.EnableRedisCache {
		client := redis.NewClient(&redis.Options{
			Addr:     config.RedisUri,
			Password: config.RedisPassword,
			DB:       config.RedisDB,
		})

		status := client.Ping(context.Background())
		redisClient = client

		log.Println("Successfully connected to Redis: ", status)
	} else {
		log.Println("Redis cache is disabled")
	}
}

func GetRedisClient(config *config.Config) *redis.Client {
	if redisClient == nil && config.EnableRedisCache {
		ConnectToRedis(config)
	}
	return redisClient
}
