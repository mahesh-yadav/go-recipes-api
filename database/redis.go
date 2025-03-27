package database

import (
	"context"

	"github.com/rs/zerolog/log"

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

		err := client.Ping(context.Background()).Err()
		if err != nil {
			log.Fatal().Err(err).Msg("Error connecting to Redis")
		}
		redisClient = client

		log.Info().Msg("Successfully connected to Redis")
	} else {
		log.Info().Msg("Redis cache is disabled")
	}
}

func GetRedisClient(config *config.Config) *redis.Client {
	if redisClient == nil && config.EnableRedisCache {
		ConnectToRedis(config)
	}
	return redisClient
}
