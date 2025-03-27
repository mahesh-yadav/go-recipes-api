package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/rs/zerolog"

	"github.com/joho/godotenv"
)

var config *Config

type Config struct {
	LogLevel                      string `env:"LOG_LEVEL" envDefault:"info"`
	LogFile                       string `env:"LOG_FILE" envDefault:"app.log"`
	LogMaxAge                     int    `env:"LOG_MAX_AGE" envDefault:"7"`
	LogMaxSizeInMB                int    `env:"LOG_MAX_SIZE_IN_MB" envDefault:"10"`
	LogCompress                   bool   `env:"LOG_COMPRESS" envDefault:"false"`
	MongoUri                      string `env:"MONGO_URI,notEmpty"`
	MongoDBName                   string `env:"MONGO_DB_NAME,notEmpty"`
	MongoServerSelectionTimeoutMS int    `env:"MONGO_SERVER_SELECTION_TIMEOUT_MS" envDefault:"5000"`
	Port                          string `env:"PORT" envDefault:"8080"`
	GinMode                       string `env:"GIN_MODE" envDefault:"debug"`
	InitializeDB                  bool   `env:"INITIALIZE_DB" envDefault:"false"`
	RedisUri                      string `env:"REDIS_URI,notEmpty"`
	RedisPassword                 string `env:"REDIS_PASSWORD"`
	RedisDB                       int    `env:"REDIS_DB" envDefault:"0"`
	EnableRedisCache              bool   `env:"ENABLE_REDIS_CACHE" envDefault:"false"`
}

func (c *Config) GetLogLevel() zerolog.Level {
	level, err := zerolog.ParseLevel(c.LogLevel)
	if err != nil {
		return zerolog.InfoLevel
	}
	return level
}

func loadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	config = &Config{}
	if err := env.Parse(config); err != nil {
		log.Fatal("error parsing config variables: ", err)
	}

	log.Println("config loaded successfully")
}

func GetConfig() *Config {
	if config == nil {
		loadConfig()
	}
	return config
}
