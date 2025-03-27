package config

import (
	"log"

	"github.com/caarlos0/env/v11"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoUri     string `env:"MONGO_URI"`
	MongoDBName  string `env:"MONGO_DB_NAME"`
	Port         string `env:"PORT" default:"8080"`
	GinMode      string `env:"GIN_MODE" default:"debug"`
	InitializeDB bool   `env:"INITIALIZE_DB" default:"false"`
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
