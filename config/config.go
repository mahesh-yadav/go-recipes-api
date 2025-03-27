package config

import (
	"log"

	"github.com/caarlos0/env/v11"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoUri         string `env:"MONGO_URI,notEmpty"`
	MongoDBName      string `env:"MONGO_DB_NAME,notEmpty"`
	Port             string `env:"PORT" envDefault:"8080"`
	GinMode          string `env:"GIN_MODE" envDefault:"debug"`
	InitializeDB     bool   `env:"INITIALIZE_DB" envDefault:"false"`
	RedisUri         string `env:"REDIS_URI,notEmpty"`
	RedisPassword    string `env:"REDIS_PASSWORD"`
	RedisDB          int    `env:"REDIS_DB" envDefault:"0"`
	EnableRedisCache bool   `env:"ENABLE_REDIS_CACHE" envDefault:"false"`
}

var config *Config

func loadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config = &Config{}
	if err := env.Parse(config); err != nil {
		log.Fatal("Error parsing environment variables: ", err)
	}

	log.Println("Config loaded successfully")
}

func GetConfig() *Config {
	if config == nil {
		loadConfig()
	}
	return config
}
